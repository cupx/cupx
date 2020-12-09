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

// alidns package implements XDns interface.
package alidns

import (
	"cupx.github.io/xdns/xdnsutil"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// AliDns implements XDns interface.
type AliDns struct {
	AK     string
	SK     string
	client *alidns.Client
}

// NewAliDns returns AliDns. 
func NewAliDns(ak, sk string) *AliDns {

	d := &AliDns{
		AK: ak,
		SK: sk,
	}
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ak, sk)
	if err != nil {
		return nil
	}
	d.client = client
	return d
}

func (d *AliDns) AddDomainRecord(t string, name string, value string) error {

	rootZone := d.GetRootZone(name)

	rr := strings.TrimSuffix(name, "."+rootZone)

	req := alidns.CreateAddDomainRecordRequest()
	req.Scheme = "https"
	req.Type = t
	req.DomainName = rootZone
	req.RR = rr
	req.Value = value
	_, err := d.client.AddDomainRecord(req)
	if err != nil {
		if strings.Contains(err.Error(), "ErrorCode: DomainRecordDuplicate") {
			return nil
		}
		return err
	}

	return nil
}

func (d *AliDns) DeleteDomainRecord(t string, name string, value string) error {
	rootZone := d.GetRootZone(name)
	rr := strings.TrimSuffix(name, "."+rootZone)
	list, err := d.AliDnsGetDomainRecordList(rootZone)
	if err != nil {
		return err
	}

	var id string

	for _, v := range list.DomainRecords.Record {
		if v.Type == t && v.RR == rr && v.Value == value {
			id = v.RecordId
		}
	}
	if id == "" {
		// don't return err if rr not found
		return nil
	}

	return d.DnsDeleteDomainRecordByID(id)

}

func (d *AliDns) DnsDeleteDomainRecordByID(id string) error {
	req := alidns.CreateDeleteDomainRecordRequest()
	req.RecordId = id
	_, err := d.client.DeleteDomainRecord(req)
	return err
}

func (d *AliDns) AliDnsGetDomainRecordList(name string) (*alidns.DescribeDomainRecordsResponse, error) {

	req := alidns.CreateDescribeDomainRecordsRequest()
	req.Scheme = "https"
	req.DomainName = name

	resp, err := d.client.DescribeDomainRecords(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (d *AliDns) GetRootZone(name string) string {

	for i := 0; ; i++ {
		subdomain := xdnsutil.TrimSubDomain(name, i)
		if subdomain == "" {
			return ""
		}

		_, err := d.AliDnsGetDomainRecordList(subdomain)
		if err == nil {
			return subdomain
		}

	}
}
