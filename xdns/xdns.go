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

// Package xdns provides an extensible dns library
package xdns

import "cupx.github.io/pkg/xdns/alidns"

// Config configures XDns when creating.
type Config struct {
	Type string
	AK   string
	SK   string
}

// XDns is common dns interface.
type XDns interface {
	// AddDomainRecord add domain record to dns server.
	AddDomainRecord(t string, name string, value string) error
	// DeleteDomainRecord delete domain record from dns server.
	DeleteDomainRecord(t string, name string, value string) error
	// DnsDeleteDomainRecordByID delete domain record from dns server by record ID. 
	DnsDeleteDomainRecordByID(id string) error
}

// NewXDns returns XDns.
func NewXDns(conf *Config) XDns {
	if conf.Type == "alidns" {
		return alidns.NewAliDns(conf.AK, conf.SK)
	}
	return nil
}
