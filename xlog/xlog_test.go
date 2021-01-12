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

package xlog

import (
	"testing"

	"cupx.github.io/pkg/xlog/xlogcore"
)

func TestNewFastXLog(t *testing.T) {
	logger, _ := NewFastXLog(xlogcore.Config{
		Level:    -1,
		Encoding: "json",
		FileName: "stdout",
	})
	defer logger.Sync()

	logger.With(xlogcore.Int64("k1", 12)).Debug("hi12")
	logger.With(xlogcore.String("k2", "v2")).Info("hi13")
	logger.Error("hi")
	logger.With(xlogcore.Int64("k2", 13)).Warn("hi14")
	logger.XLog().With("k8", "v8").Info("XLog")
	logger.XLog().Info("XLog", "Info")
	ctx := logger.With(xlogcore.String("k10", "v10")).ToCtx(nil, xlogcore.Int64("k9", 12))
	logger.FromCtx(ctx).Debug("hi from ctx")

}

func TestNewXLog(t *testing.T) {
	logger, _ := NewXLog(xlogcore.Config{
		Level:    -1,
		Encoding: "json",
		FileName: "stdout",
	})
	defer logger.Sync()

	logger.Fast().With(xlogcore.Int64("k1", 12)).Debug("hi12")
	logger.Fast().With(xlogcore.Int64("k2", 13)).Debug("hi13")
	logger.Fast().Debug("hi")

	logger.Fast().With(xlogcore.Int64("k2", 13)).Debug("hi13")

	logger.Debugf("%s%d", "hi", 12)
	logger.With("k8", "v8").Info("XLog", "k9", "v9")
	logger.Info("XLog", "k10", "v10")
	ctx := logger.ToCtx(nil, "k11", "v11")
	ctx2 := logger.ToCtx(ctx, "k14", "v14")
	logger.With("k12", "v12").FromCtx(ctx2).With("k13", "v13").Info("hi")
}

func TestNewXLogDft(t *testing.T) {
	defer Sync()

	Fast().With(xlogcore.Int64("k1", 12)).Debug("hi12")
	Fast().With(xlogcore.Int64("k2", 13)).Debug("hi13")
	Fast().Debug("hi")

	Fast().With(xlogcore.Int64("k2", 13)).Debug("hi13")

	With("k8", "v8").Info()
	Debug("XLog")
	With("k11", "v11").Warn()
	Error("hi")
	Info("msg", "msg2")
	Debugf("%v,%v", 12, "ddd")
	Infof("%v,%v", 13, "ff")
	Warnf("%v,%v", 12, "ddd")
	Errorf("%v,%v", 12, "ddd")

	Warn("msg")

	WithOptions(xlogcore.OptionWithAddCallerSkip(0)).Info("kSkip")
	AddCallerSkip(0).Info("addTest")

	ctx := ToCtx(nil, "k15", "v15")
	ctx2 := FromCtx(ctx).ToCtx(ctx, "k14", "v14")
	FromCtx(ctx2).With("k13", "v13").Info("hi")
	FromCtx(ctx2).With("k16", "v16").Fast().Info("hi")
	Fast().FromCtx(ctx2).Info("hi")
	ctx3 := ToCtx(nil)
	With("k17", "v17").FromCtx(ctx3).With("k18", "v18").Info("hi")

	//Panic("msg")
	//Fatal("msg")

	SetConfig(xlogcore.Config{
		Level:    xlogcore.DebugLevel,
		Encoding: "json",
		FileName: "",
	})
	Info("hi", "hi", 2)
	SetConfig(xlogcore.Config{
		Level:           xlogcore.DebugLevel,
		Encoding:        "json",
		FileName:        "/tmp/xlog/xlog.log",
		RotationMaxSize: 1,
	})

	for i := 0; i < 100000; i++ {
		Info("hello world.")
	}

	SetConfig(xlogcore.Config{
		Level:           xlogcore.DebugLevel,
		Encoding:        "json",
		FileName:        "/tmp/xlog/xlog.log",
		RotationMaxSize: 1,
		RetainCompress:  true,
		RetainLocalTime: true,
	})

	for i := 0; i < 100000; i++ {
		Info("hello world.")
	}

}
