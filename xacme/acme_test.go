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

package xacme

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/cupx/cupx/xacme/testdata"
	"github.com/cupx/cupx/xdns"
)

func GetStagingAcct() Client {
	data := testdata.GetTestData("../../testdata/data.test.yaml")
	dns := &xdns.Config{
		Type: data.Dns.Type,
		AK:   data.Dns.Ak,
		SK:   data.Dns.Sk,
	}
	conf := &Config{
		CA:  CaLetsencryptStaging,
		Dns: dns,
	}
	c := NewClient(
		conf,
	)
	acct := new(Account)
	acct.PemPrivateKey = data.AcmeStagingAcct.PemPrivatekey
	acct.AcctURL = data.AcmeStagingAcct.AcctURL
	acct.Contact = data.AcmeStagingAcct.Contact
	acct.TOSAgreed = data.AcmeStagingAcct.TosAgreed
	log.Println(c.SetAccount(acct))
	return c
}
func GetAcct() Client {
	data := testdata.GetTestData("./testdata/data.test.yaml")
	dns := &xdns.Config{
		Type: data.Dns.Type,
		AK:   data.Dns.Ak,
		SK:   data.Dns.Sk,
	}
	conf := &Config{
		CA:  CaLetsencrypt,
		Dns: dns,
	}
	c := NewClient(
		conf,
		WithRootCAKeyID(CaLetsencryptRootCaKeyIdIsrgRootX1),
	)
	acct := new(Account)
	acct.PemPrivateKey = data.AcmeAcct.PemPrivatekey
	acct.AcctURL = data.AcmeAcct.AcctURL
	acct.Contact = data.AcmeAcct.Contact
	acct.TOSAgreed = data.AcmeAcct.TosAgreed
	log.Println(c.SetAccount(acct))
	return c
}
func TestClient_CreatAccountWithEmail(t *testing.T) {
	dns := &xdns.Config{
		Type: "alidns",
		AK:   "",
		SK:   "",
	}
	conf := &Config{
		CA:  CaLetsencryptStaging,
		Dns: dns,
	}
	c := NewClient(
		conf,
	)

	file, err := os.OpenFile("./testdata/acct.tmp", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer file.Close()

	acctrb, err := ioutil.ReadAll(file)

	acctr := &Account{}

	json.Unmarshal(acctrb, acctr)

	log.Println(acctr)

	acct, err := c.CreateAccountWithEmail("me@archlake.net", true)
	if err != nil {
		return
	}
	acct.PrivateKey = nil

	acctb, err := json.Marshal(acct)

	log.Println(string(acctb), err)

	_, _ = file.Write(acctb)

}

func TestClient_SetAccount(t *testing.T) {
	dns := &xdns.Config{
		Type: "alidns",
		AK:   "",
		SK:   "",
	}
	conf := &Config{
		CA:  CaLetsencryptStaging,
		Dns: dns,
	}
	c := NewClient(
		conf,
	)

	file, err := os.OpenFile("./testdata/acct.tmp", os.O_RDONLY, 0777)
	if err != nil {
		return
	}
	defer file.Close()

	acctrb, err := ioutil.ReadAll(file)

	acctr := &Account{}

	json.Unmarshal(acctrb, acctr)

	log.Println(acctr, err, c)

	acct, err := c.SetAccount(acctr)

	log.Println(acct, err)

}

func TestClient_CreatAccountWithPrivateKey(t *testing.T) {
	dns := &xdns.Config{
		Type: "alidns",
		AK:   "",
		SK:   "",
	}
	conf := &Config{
		CA:  CaLetsencryptStaging,
		Dns: dns,
	}
	c := NewClient(
		conf,
	)

	file, err := os.OpenFile("./testdata/acct.tmp", os.O_RDONLY, 0777)
	if err != nil {
		return
	}
	defer file.Close()

	acctrb, err := ioutil.ReadAll(file)

	acctr := &Account{}

	json.Unmarshal(acctrb, acctr)

	log.Println(acctr, err, c)

	acct, err := c.CreateAccountWithPrivateKey(acctr)

	log.Println(acct, err)

}

func TestClient_SignCertWithDNSStaging(t *testing.T) {

	c := GetStagingAcct()

	idl := &IdlSignReq{
		Identifiers: []IdlIdentifier{
			{
				Type:  "dns",
				Value: "*.testcert1.test.xdns.cupx.net",
			},
			{
				Type:  "dns",
				Value: "testcert1.test.xdns.cupx.net",
			},
		},
		TXTCname: "testcert1.cert.issue-tls-cert.test.xdns.cupx.net",
	}

	cert, err := c.SignCertWithDNS(idl, WithRootCAKeyID(CaLetsencryptStagingRootCaKeyIdFakeLeRootX1))
	log.Println(cert, err)
	if err != nil {
		return
	}
	fmt.Println("certWithFAKE_LE_ROOT_X1:", cert.PemCertBodyWithChain)
	fmt.Println("key:", cert.PemCertPrivateKey)
	fmt.Println(cert.NotBefore)
	fmt.Println(cert.NotAfter)

	cert2, err := c.SignCertWithDNS(idl, WithRootCAKeyID(CaLetsencryptStagingRootCaKeyIdFakeLeRootX2))
	log.Println(cert2, err)
	if err != nil {
		return
	}
	fmt.Println("certWithFAKE_LE_ROOT_X2:", cert2.PemCertBodyWithChain)
	fmt.Println("key:", cert2.PemCertPrivateKey)
	fmt.Println(cert2.NotBefore)
	fmt.Println(cert2.NotAfter)
}

func TestClient_SignCertWithDNS(t *testing.T) {

	c := GetAcct()

	idl := &IdlSignReq{
		Identifiers: []IdlIdentifier{
			{
				Type:  "dns",
				Value: "testcert1.test.xdns.cupx.net",
			},
			{
				Type:  "dns",
				Value: "*.testcert1.test.xdns.cupx.net",
			},
		},
		TXTCname: "testcert1.cert.issue-tls-cert.test.xdns.cupx.net",
	}

	cert, err := c.SignCertWithDNS(idl, WithRootCAKeyID(CaLetsencryptRootCaKeyIdDstRootCaX3))
	log.Println(cert, err)
	if err != nil {
		return
	}
	fmt.Println("certWithDSTRootCA:", cert.PemCertBodyWithChain)
	fmt.Println("key:", cert.PemCertPrivateKey)
	fmt.Println(cert.NotBefore)
	fmt.Println(cert.NotAfter)

	cert2, err := c.SignCertWithDNS(idl, WithRootCAKeyID(CaLetsencryptRootCaKeyIdIsrgRootX1))
	log.Println(cert2, err)
	if err != nil {
		return
	}
	fmt.Println("certWithISRGRootCA:", cert2.PemCertBodyWithChain)
	fmt.Println("key:", cert2.PemCertPrivateKey)
	fmt.Println(cert2.NotBefore)
	fmt.Println(cert2.NotAfter)

}

func TestClient_SignCertWithDNSForDomain(t *testing.T) {

	c := GetAcct()

	idl := &IdlSignReq{
		Identifiers: []IdlIdentifier{
			{
				Type:  "dns",
				Value: "testcert1.test.xdns.cupx.net",
			},
			{
				Type:  "dns",
				Value: "*.testcert1.test.xdns.cupx.net",
			},
		},
		TXTCname: "testcert1.cert.issue-tls-cert.test.xdns.cupx.net",
	}

	cert, err := c.SignCertWithDNS(idl, WithRootCAKeyID(CaLetsencryptRootCaKeyIdIsrgRootX1))
	log.Println(cert, err)
	if err != nil {
		return
	}
	fmt.Println("certWithISRGRootCA:", cert.PemCertBodyWithChain)
	fmt.Println("key:", cert.PemCertPrivateKey)
	fmt.Println(cert.NotBefore)
	fmt.Println(cert.NotAfter)

}
