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
	"context"
	"fmt"

	"github.com/cupx/cupx/xlog/xlogcore"

	"go.uber.org/zap"
)

// XLogger wraps zap SugaredLogger.
type XLogger struct {
	logger *zap.SugaredLogger
}

// NewXLogger returns XLogger.
func NewXLogger(config xlogcore.Config) (*XLogger, error) {
	return &XLogger{
		logger: buildZapLogger(config).Sugar(),
	}, nil
}

// Sync flushes buffered logs. Users should call Sync before exiting.
func (log *XLogger) Sync() {
	log.logger.Sync()
}

// Fast converting XLog to FastXLog. The operation is quite inexpensive.
func (log *XLogger) Fast() xlogcore.FastXLog {

	return &FastXLogger{
		logger: log.logger.Desugar(),
	}
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (log *XLogger) WithOptions(opts ...xlogcore.Option) xlogcore.XLog {
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

	l.logger = l.logger.Desugar().WithOptions(zapOpts...).Sugar()
	return l
}

// AddCallerSkip increases the number of callers skipped by caller annotation.
func (log *XLogger) AddCallerSkip(n int) xlogcore.XLog {
	l := log.clone()
	l.logger = l.logger.Desugar().WithOptions(zap.AddCaller(), zap.AddCallerSkip(n)).Sugar()
	return l
}

// With adds a variadic number of fields to the logging context.
func (log *XLogger) With(kvs ...interface{}) xlogcore.XLog {

	if len(kvs) == 0 {
		return log
	}
	l := log.clone()
	l.logger = l.logger.With(kvs...)
	return l
}

// ToCtx adds the XLog, with a variadic number of fields, to ctx and
// returns the resulting context.Context.
func (log *XLogger) ToCtx(ctx context.Context, kvs ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	l := log.clone()
	l.logger = l.logger.With(kvs...)
	return context.WithValue(ctx, ctxZapXLogKey, l)
}

// FromCtx gets the XLog from the ctx.
func (log *XLogger) FromCtx(ctx context.Context) xlogcore.XLog {
	if ctx == nil {
		return log
	}

	if ctxL, ok := ctx.Value(ctxZapXLogKey).(*XLogger); ok {
		return ctxL
	}
	if ctxFl, ok := ctx.Value(ctxZapFastXLogKey).(*FastXLogger); ok {
		return ctxFl.XLog()
	}

	return log
}

// Debug construct and log a Debug message.
func (log *XLogger) Debug(args ...interface{}) {
	log.logger.Debug(log.genMsg(args...))
}

// Debugf construct and log a Debug message.
func (log *XLogger) Debugf(tpl string, args ...interface{}) {
	log.logger.Debugf(tpl, args...)
}

// Info construct and log a Info message.
func (log *XLogger) Info(args ...interface{}) {
	log.logger.Info(log.genMsg(args...))
}

// Infof construct and log a Info message.
func (log *XLogger) Infof(tpl string, args ...interface{}) {
	log.logger.Infof(tpl, args...)
}

// Warn construct and log a Warn message.
func (log *XLogger) Warn(args ...interface{}) {
	log.logger.Warn(log.genMsg(args...))
}

// Warnf construct and log a Warn message.
func (log *XLogger) Warnf(tpl string, args ...interface{}) {
	log.logger.Warnf(tpl, args...)
}

// Error construct and log a Error message.
func (log *XLogger) Error(args ...interface{}) {
	log.logger.Error(log.genMsg(args...))
}

// Errorf construct and log a Error message.
func (log *XLogger) Errorf(tpl string, args ...interface{}) {
	log.logger.Errorf(tpl, args...)
}

// Panic construct and log a Panic message, then panics.
func (log *XLogger) Panic(args ...interface{}) {
	log.logger.Panic(log.genMsg(args...))
}

// Panicf construct and log a Panic message, then panics.
func (log *XLogger) Panicf(tpl string, args ...interface{}) {
	log.logger.Panicf(tpl, args...)
}

// Fatal construct and log a Fatal message, then calls os.Exit(1).
func (log *XLogger) Fatal(args ...interface{}) {
	log.logger.Fatal(log.genMsg(args...))
}

// Fatalf construct and log a Fatal message, then calls os.Exit(1).
func (log *XLogger) Fatalf(tpl string, args ...interface{}) {
	log.logger.Fatalf(tpl, args...)
}

func (log *XLogger) clone() *XLogger {
	copy := *log
	return &copy
}

func (log *XLogger) genMsg(args ...interface{}) string {
	msg := ""
	for i, arg := range args {
		if i == 0 {
			msg += fmt.Sprintf("%v", arg)
		} else {
			msg += fmt.Sprintf(" %+v", arg)
		}
	}
	return msg
}
