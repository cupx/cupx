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
//
// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package xlogcore

import (
	"fmt"
	"math"
	"time"
)

type KeyVal struct {
	Key       string
	Type      KeyValType
	Integer   int64
	String    string
	Interface interface{}
}

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

// Skip constructs a no-op field, which is often useful when handling invalid
// inputs in other Field constructors.
func Skip() KeyVal {
	return KeyVal{Type: SkipType}
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) KeyVal {
	if err == nil {
		return Skip()
	}
	return KeyVal{Key: key, Type: ErrorType, Interface: err}
}

// nilKeyVal returns a field which will marshal explicitly as nil. See motivation
// in https://github.com/uber-go/zap/issues/753 . If we ever make breaking
// changes and add NilType and ObjectEncoder.AddNil, the
// implementation here should be changed to reflect that.
func nilKeyVal(key string) KeyVal { return Reflect(key, nil) }

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) KeyVal {
	var ival int64
	if val {
		ival = 1
	}
	return KeyVal{Key: key, Type: BoolType, Integer: ival}
}

// Boolp constructs a field that carries a *bool. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Boolp(key string, val *bool) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Bool(key, *val)
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) KeyVal {
	return KeyVal{Key: key, Type: ByteStringType, Interface: val}
}

// Complex128 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex128 to
// interface{}).
func Complex128(key string, val complex128) KeyVal {
	return KeyVal{Key: key, Type: Complex128Type, Interface: val}
}

// Complex128p constructs a field that carries a *complex128. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Complex128p(key string, val *complex128) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Complex128(key, *val)
}

// Complex64 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex64 to
// interface{}).
func Complex64(key string, val complex64) KeyVal {
	return KeyVal{Key: key, Type: Complex64Type, Interface: val}
}

// Complex64p constructs a field that carries a *complex64. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Complex64p(key string, val *complex64) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Complex64(key, *val)
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) KeyVal {
	return KeyVal{Key: key, Type: Float64Type, Integer: int64(math.Float64bits(val))}
}

// Float64p constructs a field that carries a *float64. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Float64p(key string, val *float64) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Float64(key, *val)
}

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) KeyVal {
	return KeyVal{Key: key, Type: Float32Type, Integer: int64(math.Float32bits(val))}
}

// Float32p constructs a field that carries a *float32. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Float32p(key string, val *float32) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Float32(key, *val)
}

// Int constructs a field with the given key and value.
func Int(key string, val int) KeyVal {
	return Int64(key, int64(val))
}

// Intp constructs a field that carries a *int. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Intp(key string, val *int) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Int(key, *val)
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) KeyVal {
	return KeyVal{Key: key, Type: Int64Type, Integer: val}
}

// Int64p constructs a field that carries a *int64. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Int64p(key string, val *int64) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Int64(key, *val)
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) KeyVal {
	return KeyVal{Key: key, Type: Int32Type, Integer: int64(val)}
}

// Int32p constructs a field that carries a *int32. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Int32p(key string, val *int32) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Int32(key, *val)
}

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) KeyVal {
	return KeyVal{Key: key, Type: Int16Type, Integer: int64(val)}
}

// Int16p constructs a field that carries a *int16. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Int16p(key string, val *int16) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Int16(key, *val)
}

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) KeyVal {
	return KeyVal{Key: key, Type: Int8Type, Integer: int64(val)}
}

// Int8p constructs a field that carries a *int8. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Int8p(key string, val *int8) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Int8(key, *val)
}

// String constructs a field with the given key and value.
func String(key string, val string) KeyVal {
	return KeyVal{Key: key, Type: StringType, String: val}
}

// Stringp constructs a field that carries a *string. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Stringp(key string, val *string) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return String(key, *val)
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) KeyVal {
	return Uint64(key, uint64(val))
}

// Uintp constructs a field that carries a *uint. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uintp(key string, val *uint) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uint(key, *val)
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) KeyVal {
	return KeyVal{Key: key, Type: Uint64Type, Integer: int64(val)}
}

// Uint64p constructs a field that carries a *uint64. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uint64p(key string, val *uint64) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uint64(key, *val)
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) KeyVal {
	return KeyVal{Key: key, Type: Uint32Type, Integer: int64(val)}
}

// Uint32p constructs a field that carries a *uint32. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uint32p(key string, val *uint32) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uint32(key, *val)
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) KeyVal {
	return KeyVal{Key: key, Type: Uint16Type, Integer: int64(val)}
}

// Uint16p constructs a field that carries a *uint16. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uint16p(key string, val *uint16) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uint16(key, *val)
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) KeyVal {
	return KeyVal{Key: key, Type: Uint8Type, Integer: int64(val)}
}

// Uint8p constructs a field that carries a *uint8. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uint8p(key string, val *uint8) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uint8(key, *val)
}

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, val uintptr) KeyVal {
	return KeyVal{Key: key, Type: UintptrType, Integer: int64(val)}
}

// Uintptrp constructs a field that carries a *uintptr. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Uintptrp(key string, val *uintptr) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Uintptr(key, *val)
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) KeyVal {
	return KeyVal{Key: key, Type: ReflectType, Interface: val}
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) KeyVal {
	return KeyVal{Key: key, Type: NamespaceType}
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) KeyVal {
	return KeyVal{Key: key, Type: StringerType, Interface: val}
}

// Time constructs a KeyVal with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) KeyVal {
	if val.Before(_minTimeInt64) || val.After(_maxTimeInt64) {
		return KeyVal{Key: key, Type: TimeFullType, Interface: val}
	}
	return KeyVal{Key: key, Type: TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}

// Timep constructs a field that carries a *time.Time. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Timep(key string, val *time.Time) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Time(key, *val)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) KeyVal {
	return KeyVal{Key: key, Type: DurationType, Integer: int64(val)}
}

// Durationp constructs a field that carries a *time.Duration. The returned KeyVal will safely
// and explicitly represent `nil` when appropriate.
func Durationp(key string, val *time.Duration) KeyVal {
	if val == nil {
		return nilKeyVal(key)
	}
	return Duration(key, *val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) KeyVal {
	switch val := value.(type) {
	case bool:
		return Bool(key, val)
	case *bool:
		return Boolp(key, val)
	case complex128:
		return Complex128(key, val)
	case *complex128:
		return Complex128p(key, val)
	case complex64:
		return Complex64(key, val)
	case *complex64:
		return Complex64p(key, val)
	case float64:
		return Float64(key, val)
	case *float64:
		return Float64p(key, val)
	case float32:
		return Float32(key, val)
	case *float32:
		return Float32p(key, val)
	case int:
		return Int(key, val)
	case *int:
		return Intp(key, val)
	case int64:
		return Int64(key, val)
	case *int64:
		return Int64p(key, val)
	case int32:
		return Int32(key, val)
	case *int32:
		return Int32p(key, val)
	case int16:
		return Int16(key, val)
	case *int16:
		return Int16p(key, val)
	case int8:
		return Int8(key, val)
	case *int8:
		return Int8p(key, val)
	case string:
		return String(key, val)
	case *string:
		return Stringp(key, val)
	case uint:
		return Uint(key, val)
	case *uint:
		return Uintp(key, val)
	case uint64:
		return Uint64(key, val)
	case *uint64:
		return Uint64p(key, val)
	case uint32:
		return Uint32(key, val)
	case *uint32:
		return Uint32p(key, val)
	case uint16:
		return Uint16(key, val)
	case *uint16:
		return Uint16p(key, val)
	case uint8:
		return Uint8(key, val)
	case *uint8:
		return Uint8p(key, val)
	case uintptr:
		return Uintptr(key, val)
	case *uintptr:
		return Uintptrp(key, val)
	case time.Time:
		return Time(key, val)
	case *time.Time:
		return Timep(key, val)
	case time.Duration:
		return Duration(key, val)
	case *time.Duration:
		return Durationp(key, val)
	case error:
		return NamedError(key, val)
	case fmt.Stringer:
		return Stringer(key, val)
	default:
		return Reflect(key, val)
	}
}
