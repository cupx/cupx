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

package testdata

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type TestData struct {
	Dns struct {
		Type string `yaml:"type"`
		Ak   string `yaml:"ak"`
		Sk   string `yaml:"sk"`
	} `yaml:"dns"`
}

func GetTestData(path string) *TestData {
	file, err := os.OpenFile(path, os.O_RDWR, 777)
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	t := new(TestData)
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return nil
	}
	return t
}
