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
	"context"
	
	"github.com/cupx/cupx/xlog/xlogcore"
)

var dft xlogcore.XLog

func init() {
	dft, _ = NewXLog(xlogcore.Config{
		Level:    xlogcore.DebugLevel,
		Encoding: "console",
		FileName: "stdout",
	})
}

// SetConfig reconfigures the dft logger.
func SetConfig(cfg xlogcore.Config) (err error) {
	dft, err = NewXLog(cfg)
	return err
}

// Sync flushes buffered logs. Users should call Sync before exiting.
func Sync() {
	dft.Sync()
}

// Fast converting XLog to FastXLog. The operation is quite inexpensive.
func Fast() xlogcore.FastXLog {
	return dft.Fast()
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...xlogcore.Option) xlogcore.XLog {
	return dft.WithOptions(opts...)
}

// AddCallerSkip increases the number of callers skipped by caller annotation.
func AddCallerSkip(n int) xlogcore.XLog {
	return dft.AddCallerSkip(n)
}

// With adds a variadic number of fields to the logging context.
func With(kvs ...interface{}) xlogcore.XLog {
	return dft.With(kvs...)
}

// Withc adds a variadic number of fields to the context.Context and
// returns the resulting ctx.
func Withc(ctx context.Context, kvs ...interface{}) context.Context {
	return dft.Withc(ctx, kvs...)
}

// Ctx adds a variadic number of fields to the logging context from the ctx.
func Ctx(ctx context.Context) xlogcore.XLog {
	return dft.Ctx(ctx)
}

// Debug construct and log a Debug message.
func Debug(args ...interface{}) {
	dft.AddCallerSkip(1).Debug(args...)
}

// Debugf construct and log a Debug message.
func Debugf(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Debugf(tpl, args...)
}

// Info construct and log a Info message.
func Info(args ...interface{}) {
	dft.AddCallerSkip(1).Info(args...)
}

// Infof construct and log a Info message.
func Infof(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Infof(tpl, args...)
}

// Warn construct and log a Warn message.
func Warn(args ...interface{}) {
	dft.AddCallerSkip(1).Warn(args...)
}

// Warnf construct and log a Warn message.
func Warnf(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Warnf(tpl, args...)
}

// Error construct and log a Error message.
func Error(args ...interface{}) {
	dft.AddCallerSkip(1).Error(args...)
}

// Errorf construct and log a Error message.
func Errorf(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Errorf(tpl, args...)
}

// Panic construct and log a Panic message, then panics.
func Panic(args ...interface{}) {
	dft.AddCallerSkip(1).Panic(args...)
}

// Panicf construct and log a Panic message.
func Panicf(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Panicf(tpl, args...)
}

// Fatal construct and log a Fatal message, then calls os.Exit(1).
func Fatal(args ...interface{}) {
	dft.AddCallerSkip(1).Fatal(args...)
}

// Fatalf construct and log a Fatal message.
func Fatalf(tpl string, args ...interface{}) {
	dft.AddCallerSkip(1).Fatalf(tpl, args...)
}
