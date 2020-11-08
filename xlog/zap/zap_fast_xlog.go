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
	"github.com/cupx/cupx/xlog/xlogcore"
	"go.uber.org/zap"
)

// XLogger wraps zap Logger.
type FastXLogger struct {
	logger *zap.Logger
}

// NewFastXLogger returns FastXLogger.
func NewFastXLogger(config xlogcore.Config) (*FastXLogger, error) {
	return &FastXLogger{
		logger: buildZapLogger(config),
	}, nil
}

// Sync flushes buffered logs. Users should call Sync before exiting.
func (log *FastXLogger) Sync() {
	_ = log.logger.Sync()
}

// XLog converting FastXLog to XLog. The operation is quite inexpensive.
func (log *FastXLogger) XLog() xlogcore.XLog {
	return &XLogger{
		logger: log.logger.Sugar(),
	}
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (log *FastXLogger) WithOptions(opts ...xlogcore.Option) xlogcore.FastXLog {
	l := log.clone()

	o := &xlogcore.Options{}

	for _, opt := range opts {
		opt(o)
	}
	zapOpts := make([]zap.Option, 0)

	if o.AddCallerSkip != nil {
		zapOpts = append(zapOpts, zap.AddCaller())
		zapOpts = append(zapOpts, zap.AddCallerSkip(*o.AddCallerSkip))
	}

	l.logger = log.logger.WithOptions(zapOpts...)
	return l
}

// AddCallerSkip increases the number of callers skipped by caller annotation.	
func (log *FastXLogger) AddCallerSkip(n int) xlogcore.FastXLog {
	l := log.clone()
	l.logger = log.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(n))
	return l
}

// With adds a variadic number of fields to the logging context.	
func (log *FastXLogger) With(kvs ...xlogcore.KeyVal) xlogcore.FastXLog {

	if len(kvs) == 0 {
		return log
	}
	l := log.clone()
	l.logger = log.logger.With(log.keyValToZapField(kvs...)...)
	return l
}

// Debug construct and log a Debug message.
func (log *FastXLogger) Debug(msg string) {
	log.logger.Debug(msg)
}

// Info construct and log a Info message.
func (log *FastXLogger) Info(msg string) {
	log.logger.Info(msg)
}

// Warn construct and log a Warn message.
func (log *FastXLogger) Warn(msg string) {
	log.logger.Warn(msg)
}

// Error construct and log a Error message.
func (log *FastXLogger) Error(msg string) {
	log.logger.Error(msg)
}

// Panic construct and log a Panic message, then panics.
func (log *FastXLogger) Panic(msg string) {
	log.logger.Panic(msg)
}

// Fatal construct and log a Fatal message, then calls os.Exit(1).
func (log *FastXLogger) Fatal(msg string) {
	log.logger.Fatal(msg)
}

func (log *FastXLogger) clone() *FastXLogger {
	copy := *log
	return &copy
}

func (log *FastXLogger) keyValToZapField(kvs ...xlogcore.KeyVal) []zap.Field {

	var fields []zap.Field
	for _, kv := range kvs {
		switch kv.Type {
		case xlogcore.BoolType:
			fields = append(fields, zap.Bool(kv.Key, func() bool {
				if kv.Integer == 1 {
					return true
				}
				return false
			}()))
		case xlogcore.StringType:
			fields = append(fields, zap.String(kv.Key, kv.String))
		case xlogcore.Int64Type:
			fields = append(fields, zap.Int64(kv.Key, kv.Integer))
		case xlogcore.ByteStringType:
			fields = append(fields, zap.ByteString(kv.Key, func() []byte {
				v, ok := kv.Interface.([]byte)
				if ok {
					return v
				}
				return nil
			}()))
		default:
			fields = append(fields, zap.Any(kv.Key, kv.Interface))
		}
	}

	return fields
}
