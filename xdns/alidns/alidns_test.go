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

package alidns

import (
	"log"
	"testing"

	"cupx.github.io/pkg/xdns/alidns/testdata"
)

func newClient() *AliDns {
	data := testdata.GetTestData("./testdata/data.test.yaml")
	return NewAliDns(
		data.Dns.Ak,
		data.Dns.Sk,
	)
}

func TestAliDns_AliDnsGetDomainRecordList(t *testing.T) {
	d := newClient()

	log.Println(d.AliDnsGetDomainRecordList("test.xdns.cupx.net"))
}

func TestAliDns_GetRootZone(t *testing.T) {
	d := newClient()
	log.Println(d.GetRootZone("t1.test.test.xdns.cupx.net"))
}

func TestAliDns_AddDomainRecord(t *testing.T) {
	d := newClient()

	log.Println(d.AddDomainRecord("TXT", "q.w.dc.x.test.xdns.cupx.net", "test"))
}

func TestAliDns_DeleteDomainRecord(t *testing.T) {
	d := newClient()
	log.Println(d.DeleteDomainRecord("TXT", "q.w.dc.x.test.xdns.cupx.net", "test"))
}
