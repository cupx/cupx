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
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"gopkg.in/square/go-jose.v2"
)

func GetJWKThumbprintWithBase64url(key interface{}) (string, error) {
	j := jose.JSONWebKey{
		Key: key,
	}

	jd, err := j.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", nil
	}

	return base64.RawURLEncoding.EncodeToString(jd), nil

}

func Sha256WithBase64url(b []byte) (d string) {
	h := sha256.New()

	h.Write(b)

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

type HTTPHeaderLink struct {
	URL    string
	Rel    string
	Params map[string]string
}

func GetHTTPHeaderLink(ss []string) []HTTPHeaderLink {
	var links []HTTPHeaderLink
	for _, chunk := range ss {
		link := HTTPHeaderLink{Params: make(map[string]string)}
		for _, piece := range strings.Split(chunk, ";") {

			piece = strings.Trim(piece, " ")
			if piece == "" {
				continue
			}

			if piece[0] == '<' && piece[len(piece)-1] == '>' {
				link.URL = strings.Trim(piece, "<>")
				continue
			}

			key, val := "", ""
			parts := strings.SplitN(piece, "=", 2)
			if len(parts) == 1 {
				key = parts[0]
			} else if len(parts) == 2 {
				key = parts[0]
				val = strings.Trim(parts[1], "\"")
			}

			if key == "" {
				continue
			}

			if strings.ToLower(key) == "rel" {
				link.Rel = val
			} else {
				link.Params[key] = val
			}
		}
		if link.URL != "" {
			links = append(links, link)
		}
	}
	return links
}

func FmtX509KeyID(id []byte) string {
	s := ""
	for k, v := range id {
		if k < len(id)-1 {
			s += fmt.Sprintf("%02X:", v)
		} else {
			s += fmt.Sprintf("%02X", v)
		}

	}
	return s
}
