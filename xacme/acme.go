// Copyright 2020 The CupX Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// xacme package implement part of rfc8555.
// https://tools.ietf.org/html/rfc8555
package xacme

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cupx/cupx/xdns"
	"gopkg.in/square/go-jose.v2"
)

const (
	CaLetsencrypt                               = "letsencrypt"
	CaLetsencryptStaging                        = "letsencrypt_staging"
	CaLetsencryptRootCaKeyIdIsrgRootX1          = "79:B4:59:E6:7B:B6:E5:E4:01:73:80:08:88:C8:1A:58:F6:E9:9B:6E"
	CaLetsencryptRootCaKeyIdDstRootCaX3         = "C4:A7:B1:A4:7B:2C:71:FA:DB:E1:4B:90:75:FF:C4:15:60:85:89:10"
	CaLetsencryptStagingRootCaKeyIdFakeLeRootX1 = "C1:26:74:A4:8A:44:A0:E6:FA:20:28:D8:5C:23:9A:45:88:18:79:E0"
	CaLetsencryptStagingRootCaKeyIdFakeLeRootX2 = "1B:FB:1C:F0:31:7D:03:2B:DA:0A:9B:AF:78:A6:F6:99:91:19:9C:B2"
)

// CaMeta contains the Directory URL.
type CaMeta struct {
	DirURL      string
	NewAcctURL  string
	NewOrderURL string
}

var caAcmeDirMap = map[string]string{
	CaLetsencrypt:        "https://acme-v02.api.letsencrypt.org/directory",
	CaLetsencryptStaging: "https://acme-staging-v02.api.letsencrypt.org/directory",
}

// Client is the acme client interface.
type Client interface {
	// CreateAccountWithEmail create acme account with email.
	CreateAccountWithEmail(email string, TOSAgreed bool) (*Account, error)
	// SetAccount set Account for acme client.
	SetAccount(acct *Account) (*Account, error)
	// CreateAccountWithPrivateKey create acme account with private key. 
	CreateAccountWithPrivateKey(acct *Account) (*Account, error)
	// SignCertWithDNS sign certificate with dns-01 Challenge. 
	SignCertWithDNS(sr *IdlSignReq, opts ...Option) (*CertInfo, error)
}

// Option configures option.
type Option func(opt *option)

type option struct {
	RootCAKeyID string
}

// WithRootCAKeyID chooses which Root CA to use.
func WithRootCAKeyID(id string) Option {
	return func(opt *option) {
		opt.RootCAKeyID = id
	}
}

type client struct {
	ca         string
	acct       *Account
	httpClient http.Client
	nonce      string
	dns        *xdns.Config
	caMeta     *CaMeta
	postMu     *sync.Mutex
	opt        option
}

// Config configures a Client when creating.
type Config struct {
	CA  string
	Dns *xdns.Config
}

// NewClient return a acme client.
func NewClient(conf *Config, opts ...Option) Client {

	d, ok := caAcmeDirMap[conf.CA]
	if !ok {
		return nil
	}
	req, err := http.NewRequest(http.MethodGet, d, nil)
	if err != nil {
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	idl := &IdlRespDir{}
	err = json.Unmarshal(respb, idl)
	if err != nil {
		return nil
	}

	return &client{
		ca:         conf.CA,
		httpClient: http.Client{},
		dns:        conf.Dns,
		caMeta: &CaMeta{
			NewAcctURL:  idl.NewAccount,
			NewOrderURL: idl.NewOrder,
		},
		postMu: new(sync.Mutex),
		opt: func() option {
			o := new(option)
			for _, opt := range opts {
				opt(o)
			}
			return *o
		}(),
	}
}

// Account contains acme account data.
type Account struct {
	Contact       []string
	TOSAgreed     bool
	AcctURL       string
	PrivateKey    *ecdsa.PrivateKey
	PemPrivateKey string
}

// Account contains signed cert info.
type CertInfo struct {
	SignatureAlgorithm   string
	PemCertPrivateKey    string
	PemCertChain         string
	PemCertBody          string
	PemCertBodyWithChain string
	NotBefore            string
	NotAfter             string
	RootCAKeyID          string
}

func (c *client) SetAccount(acct *Account) (*Account, error) {

	if acct.PemPrivateKey != "" {
		block, _ := pem.Decode([]byte(acct.PemPrivateKey))
		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		acct.PrivateKey = privateKey
	}

	c.acct = acct

	return c.acct, nil
}

func (c *client) CreateAccountWithEmail(email string, TOSAgreed bool) (*Account, error) {
	c.acct = &Account{
		Contact:   []string{email},
		TOSAgreed: TOSAgreed,
	}

	err := c.newAccount()
	if err != nil {
		return nil, err
	}

	return c.acct, nil
}

func (c *client) CreateAccountWithPrivateKey(acct *Account) (*Account, error) {

	_, err := c.SetAccount(acct)
	if err != nil {
		return nil, err
	}

	err = c.newAccount()
	if err != nil {
		return nil, err
	}

	return c.acct, nil
}

func (c *client) SignCertWithDNS(sr *IdlSignReq, opts ...Option) (*CertInfo, error) {

	nc := c.clone()
	for _, opt := range opts {
		opt(&nc.opt)
	}
	// new order.
	o := &IdlReqNewOrderPayload{
		Identifiers: sr.Identifiers,
	}
	oResp, err := nc.newOrder(o)
	if err != nil {
		return nil, err
	}

	// validate identifier with dns. 
	err = nc.validateIdentifierWithDNS(oResp.Authorizations, sr.TXTCname)
	if err != nil {
		return nil, err
	}

	// create csp.
	csr, pri, err := nc.getCsrUseSHA256WithRSA(sr)
	pemPri := string(pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: func() []byte {
			priB := x509.MarshalPKCS1PrivateKey(pri)
			if err != nil {
				return nil
			}
			return priB
		}(),
	}))

	// request certificate.
	fRespB, _, err := nc.acmePost(oResp.Finalize, fmt.Sprintf("{\"csr\":\"%s\"}", csr), true)

	// download certificate.
	fResp := &IdlRespFinalize{}
	err = json.Unmarshal(fRespB, fResp)
	if err != nil {
		return nil, err
	}
	if fResp.Status == "valid" {
		return nc.getCertFromURL(fResp.Certificate, pemPri)
	}

	return nil, errors.New("order status not valid")
}

func (c *client) newAccount() error {

	if c.acct.PrivateKey == nil {

		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		c.acct.PrivateKey = privateKey
	}

	c.acct.PemPrivateKey = string(pem.EncodeToMemory(&pem.Block{
		Type: "EC PRIVATE KEY",
		Bytes: func() []byte {
			prib, err := x509.MarshalECPrivateKey(c.acct.PrivateKey)
			if err != nil {
				return nil
			}
			return prib
		}(),
	}))

	var contact []string
	for _, v := range c.acct.Contact {
		contact = append(contact, "mailto:"+v)
	}

	payload := IdlReqNewAccountPayload{
		TermsOfServiceAgreed: true,
		Contact:              contact,
	}

	_, resp, err := c.acmePost(c.caMeta.NewAcctURL, payload, false)

	if err != nil {
		return err
	}

	c.acct.AcctURL = resp.Header.Get("Location")

	return nil
}

func (c *client) clone() *client {
	nc := *c
	nc.postMu = new(sync.Mutex)
	nc.nonce = ""
	return &nc
}
func (c *client) getCertFromURL(url string, pemPri string) (*CertInfo, error) {

	var certPems [][]byte

	DcRespB, DcResp, err := c.acmePost(url, "", true)
	if err != nil {
		return nil, err
	}
	certPems = append(certPems, DcRespB)

	links := GetHTTPHeaderLink(DcResp.Header.Values("Link"))

	for _, link := range links {
		if link.Rel == "alternate" {
			DcRespB, _, err := c.acmePost(link.URL, "", true)
			if err != nil {
				continue
			}
			certPems = append(certPems, DcRespB)
		}
	}

	var certInfos []*CertInfo
	for i, certPem := range certPems {
		var cert509s []*x509.Certificate
		var certBlocks []*pem.Block
		for len(certPem) > 0 {
			block, restPem := pem.Decode(certPem)
			certPem = restPem
			cert509, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			cert509s = append(cert509s, cert509)
			certBlocks = append(certBlocks, block)
		}

		certInfo := new(CertInfo)
		if cert509s == nil || certBlocks == nil || len(cert509s) < 1 || len(certBlocks) < 1 {
			continue
		}

		certInfo.RootCAKeyID = FmtX509KeyID(cert509s[len(cert509s)-1].AuthorityKeyId)
		if c.opt.RootCAKeyID != "" {
			if certInfo.RootCAKeyID != c.opt.RootCAKeyID {
				continue
			}
		}

		certInfo.NotBefore = cert509s[0].NotBefore.UTC().Format(time.RFC3339)
		certInfo.NotAfter = cert509s[0].NotAfter.UTC().Format(time.RFC3339)
		certInfo.PemCertBodyWithChain = string(certPems[i])
		certInfo.PemCertBody = string(pem.EncodeToMemory(certBlocks[0]))
		for i := 1; i < len(certBlocks); i++ {
			certInfo.PemCertChain += string(pem.EncodeToMemory(certBlocks[i]))
		}
		certInfo.SignatureAlgorithm = cert509s[0].SignatureAlgorithm.String()
		certInfo.PemCertPrivateKey = pemPri
		certInfos = append(certInfos, certInfo)
	}
	if len(certInfos) >= 1 {
		return certInfos[0], nil
	}

	return nil, errors.New("failed to get certInfo")
}

func (c *client) validateIdentifierWithDNS(authzs []string, cname string) error {

	var wg sync.WaitGroup
	for _, authz := range authzs {
		wg.Add(1)
		go func(authz string) {
			defer wg.Done()
			_ = c.dns01Challenge(authz, cname)
		}(authz)
	}
	wg.Wait()
	for _, authz := range authzs {
		darResp, err := c.downloadAuthorizationResources(authz)
		if err != nil {
			return err
		}
		if darResp.Status != "valid" {
			return fmt.Errorf("authorization err, status : %s", darResp.Status)
		}
	}
	return nil
}

func (c *client) dns01Challenge(authz string, cname string) error {
	darResp, err := c.downloadAuthorizationResources(authz)
	if err != nil {
		return err
	}
	if darResp.Status != "pending" {
		return nil
	}
	for _, challenge := range darResp.Challenges {
		if challenge.Type == "dns-01" {
			var name string
			if cname != "" {
				name = cname
			} else {
				name = "_acme-challenge." + darResp.Identifier.Value
			}

			tp, err := GetJWKThumbprintWithBase64url(c.acct.PrivateKey.Public())
			if err != nil {
				return err
			}

			keyAuthDigest := Sha256WithBase64url([]byte(challenge.Token + "." + tp))
			dns := xdns.NewXDns(c.dns)
			err = dns.AddDomainRecord("TXT", name, keyAuthDigest)
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 30)
			_, _, err = c.acmePost(challenge.URL, "{}", true)

			if err != nil {
				return err
			}

			checkCount := 0
			for {
				checkCount++
				darResp, err := c.downloadAuthorizationResources(authz)
				if err != nil {
					return err
				}
				if darResp.Status == "pending" {
					time.Sleep(time.Second * 5)
					if checkCount < 20 {
						continue
					}
				}
				break
			}

			err = dns.DeleteDomainRecord("TXT", name, keyAuthDigest)
			if err != nil {
				return err
			}

			break
		}

	}
	return nil
}

func (c *client) getCsrUseSHA256WithRSA(sr *IdlSignReq) (string, *rsa.PrivateKey, error) {

	crt := &x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName: sr.Identifiers[0].Value,
		},
		DNSNames: func() []string {

			var names []string
			for _, identifier := range sr.Identifiers {
				names = append(names, identifier.Value)
			}

			return names
		}(),
	}
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", nil, err
	}

	csrByte, err := x509.CreateCertificateRequest(rand.Reader, crt, priKey)
	if err != nil {
		return "", nil, err
	}

	csr := base64.RawURLEncoding.EncodeToString(csrByte)

	return csr, priKey, nil
}

func (c *client) newOrder(p *IdlReqNewOrderPayload) (*IdlRespNewOrder, error) {

	resps, _, err := c.acmePost(c.caMeta.NewOrderURL, p, true)
	if err != nil {
		return nil, err
	}

	res := &IdlRespNewOrder{}

	err = json.Unmarshal(resps, res)
	if err != nil {
		return nil, err
	}

	if res.Status == "pending" || res.Status == "ready" {
		return res, nil
	}

	return nil, fmt.Errorf("err New order")

}

func (c *client) downloadAuthorizationResources(url string) (*IdlRespDownLoadAuthorizationResources, error) {

	respB, _, err := c.acmePost(url, nil, true)
	if err != nil {
		return nil, err
	}

	idl := &IdlRespDownLoadAuthorizationResources{}

	err = json.Unmarshal(respB, idl)
	if err != nil {
		return nil, err
	}

	return idl, nil

}

func (c *client) signPayloadWithES256(url string, p interface{}, hasKid bool) (string, error) {

	s, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.ES256,
			Key:       c.acct.PrivateKey,
		},
		&jose.SignerOptions{
			EmbedJWK: func() bool {
				if hasKid {
					return false
				}
				return true
			}(),
			ExtraHeaders: func() map[jose.HeaderKey]interface{} {

				h := make(map[jose.HeaderKey]interface{}, 0)
				h["nonce"] = c.nonce
				h["url"] = url
				if hasKid {
					h["kid"] = c.acct.AcctURL
				}
				return h
			}(),
		},
	)

	if err != nil {
		return "", err
	}

	var jp []byte

	if ps, ok := p.(string); ok {
		jp = []byte(ps)
	} else {
		if p != nil {
			jp, err = json.Marshal(p)

			if err != nil {
				return "", err
			}
		} else {
			jp = []byte("")
		}
	}

	jws, err := s.Sign(jp)

	if err != nil {
		return "", err
	}

	return jws.FullSerialize(), nil

}

func (c *client) acmePost(url string, p interface{}, hasKid bool) (respb []byte, resp *http.Response, err error) {
	c.postMu.Lock()
	defer c.postMu.Unlock()

	SignAndPostF := func(url string, p interface{}, hasKid bool) (respb []byte, resp *http.Response, err error) {
		sg, _ := c.signPayloadWithES256(url, p, hasKid)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(sg)))

		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/jose+json")

		resp, err = c.httpClient.Do(req)

		if err != nil {
			return
		}

		respb, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return
		}

		if nonce := resp.Header.Get("Replay-Nonce"); nonce != "" {
			c.nonce = nonce
		}
		return
	}

	respb, resp, err = SignAndPostF(url, p, hasKid)

	if resp != nil && resp.StatusCode == http.StatusBadRequest {

		respe := &IdlRespErr{}

		_ = json.Unmarshal(respb, respe)

		if respe.Type == "urn:ietf:params:acme:error:badNonce" {
			respb, resp, err = SignAndPostF(url, p, hasKid)
		}
	}

	if resp != nil && resp.StatusCode >= 400 {
		respe := &IdlRespErr{}

		_ = json.Unmarshal(respb, respe)

		err = errors.New(respe.Type + " " + respe.Detail)
	}

	return

}
