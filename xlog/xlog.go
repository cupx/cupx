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

// Package xlog provides an extensible acme library
package xlog

import (
	"fmt"
	
	"cupx.github.io/xlog/xlogcore"
	"cupx.github.io/xlog/zap"
)

// NewFastXLog returns FastXLog.
func NewFastXLog(config xlogcore.Config) (xlogcore.FastXLog, error) {
	config = buildConfig(config)
	if config.Driver == "zap" {
		return zap.NewFastXLogger(config)
	}
	return nil, fmt.Errorf("the driver %s is not supported", config.Driver)
}

// NewXLog returns XLog.
func NewXLog(config xlogcore.Config) (xlogcore.XLog, error) {
	config = buildConfig(config)
	if config.Driver == "zap" {
		return zap.NewXLogger(config)
	}
	return nil, fmt.Errorf("the driver %s is not supported", config.Driver)
}

func buildConfig(config xlogcore.Config) xlogcore.Config {
	if config.Driver == "" {
		config.Driver = "zap"
	}

	if config.FileName == "" {
		config.FileName = "stdout"
	}

	return config
}
