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

package zap

import (
	"os"

	"github.com/cupx/cupx/xlog/xlogcore"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildZapLogger(config xlogcore.Config) *zap.Logger {

	if config.Level >= xlogcore.PanicLevel {
		config.Level++
	}

	// according to https://github.com/uber-go/zap/blob/master/FAQ.md#does-zap-support-log-rotation
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	var w zapcore.WriteSyncer
	if config.FileName == "stdout" {
		w = zapcore.AddSync(os.Stdout)
	} else {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FileName,
			MaxSize:    config.RotationMaxSize, // megabytes
			MaxBackups: config.RetainMaxBackups,
			MaxAge:     config.RetainMaxAge, // days
			LocalTime:  config.RetainLocalTime,
			Compress:   config.RetainCompress,
		})
	}

	encCfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var enc zapcore.Encoder
	if config.Encoding == "json" {
		enc = zapcore.NewJSONEncoder(encCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encCfg)
	}

	core := zapcore.NewCore(
		enc,
		w,
		zapcore.Level(config.Level),
	)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
