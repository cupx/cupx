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

package xlogcore

import "context"

// FastXLog provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
// Xlog can be converted to a FastXLog with its Fast method.
type FastXLog interface {
	// Sync flushes buffered logs. Users should call Sync before exiting.
	Sync()

	// XLog converting FastXLog to XLog. The operation is quite inexpensive.	
	XLog() XLog
	// WithOptions clones the current Logger, applies the supplied Options, and
	// returns the resulting Logger. It's safe to use concurrently.
	WithOptions(opts ...Option) FastXLog
	// AddCallerSkip increases the number of callers skipped by caller annotation.	
	AddCallerSkip(n int) FastXLog

	// With adds a variadic number of fields to the logging context.	
	With(kvs ...KeyVal) FastXLog

	// ToCtx adds the FastXLog, with a variadic number of fields, to ctx and
	// returns the resulting context.Context.
	ToCtx(ctx context.Context, kvs ...KeyVal) context.Context
	// FromCtx gets the FastXLog from the ctx.
	FromCtx(ctx context.Context) FastXLog

	// Debug construct and log a Debug message.
	Debug(msg string)
	// Info construct and log a Info message.
	Info(msg string)
	// Warn construct and log a Warn message.
	Warn(msg string)
	// Error construct and log a Error message.
	Error(msg string)
	// Panic construct and log a Panic message, then panics.
	Panic(msg string)
	// Fatal construct and log a Fatal message, then calls os.Exit(1).
	Fatal(msg string)
}

// Xlog provide a more ergonomic, but slightly slower, API.
// FastXLog can be converted to a Xlog with its Xlog method.
type XLog interface {
	// Sync flushes buffered logs. Users should call Sync before exiting.
	Sync()

	// Fast converting XLog to FastXLog. The operation is quite inexpensive.
	Fast() FastXLog
	// WithOptions clones the current Logger, applies the supplied Options, and
	// returns the resulting Logger. It's safe to use concurrently.
	WithOptions(opts ...Option) XLog
	// AddCallerSkip increases the number of callers skipped by caller annotation.
	AddCallerSkip(n int) XLog

	// With adds a variadic number of fields to the logging context.
	With(kvs ...interface{}) XLog

	// ToCtx adds the XLog, with a variadic number of fields, to ctx and
	// returns the resulting context.Context.
	ToCtx(ctx context.Context, kvs ...interface{}) context.Context
	// FromCtx gets the XLog from the ctx.
	FromCtx(ctx context.Context) XLog

	// Debug construct and log a Debug message.
	Debug(args ...interface{})
	// Debugf construct and log a Debug message.
	Debugf(tpl string, args ...interface{})
	// Info construct and log a Info message.
	Info(args ...interface{})
	// Infof construct and log a Info message.
	Infof(tpl string, args ...interface{})
	// Warn construct and log a Warn message.
	Warn(args ...interface{})
	// Warnf construct and log a Warn message.
	Warnf(tpl string, args ...interface{})
	// Error construct and log a Error message.
	Error(args ...interface{})
	// Errorf construct and log a Error message.
	Errorf(tpl string, args ...interface{})
	// Panic construct and log a Panic message, then panics.
	Panic(args ...interface{})
	// Panicf construct and log a Panic message, then panics.
	Panicf(tpl string, args ...interface{})
	// Fatal construct and log a Fatal message, then calls os.Exit(1).
	Fatal(args ...interface{})
	// Fatalf construct and log a Fatal message, then calls os.Exit(1).
	Fatalf(tpl string, args ...interface{})
}
