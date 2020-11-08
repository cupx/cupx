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

type IdlRespErr struct {
	Type   string
	Detail string
}

type IdlSignReq struct {
	Identifiers []IdlIdentifier
	TXTCname    string
}
type IdlRespDir struct {
	KeyChange string `json:"keyChange"`
	Meta      struct {
		CaaIdentities  []string `json:"caaIdentities"`
		TermsOfService string   `json:"termsOfService"`
		Website        string
	}
	NewAccount string `json:"newAccount"`
	NewNonce   string `json:"newNonce"`
	NewOrder   string `json:"newOrder"`
	RevokeCert string `json:"revokeCert"`
}

type IdlReqNewAccountPayload struct {
	TermsOfServiceAgreed bool `json:"termsOfServiceAgreed"`
	Contact              []string
}

type IdlRespNewAccount struct {
	Status  string
	Contact []string
	Order   string
}

type IdlIdentifier struct {
	Type  string
	Value string
}

type IdlReqNewOrderPayload struct {
	Identifiers []IdlIdentifier
	NotBefore   string `json:"NotBefore"`
	NotAfter    string `json:"NotAfter"`
}

type IdlRespNewOrder struct {
	Status         string
	Expires        string
	NotBefore      string `json:"NotBefore"`
	NotAfter       string `json:"NotAfter"`
	Identifiers    []IdlIdentifier
	Authorizations []string
	Finalize       string
}

type IdlChallenge struct {
	Type  string
	URL   string `json:"url"`
	Token string
}
type IdlRespDownLoadAuthorizationResources struct {
	Status     string
	Expires    string
	Identifier IdlIdentifier
	Challenges []IdlChallenge
}

type IdlRespFinalize struct {
	Status         string
	Expires        string
	NotBefore      string `json:"NotBefore"`
	NotAfter       string `json:"NotAfter"`
	Identifiers    []IdlIdentifier
	Authorizations []string
	Finalize       string
	Certificate    string
}
