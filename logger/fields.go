package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Field = zap.Field

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Bools(key string, val []bool) Field {
	return zap.Bools(key, val)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func ByteStrings(key string, val [][]byte) Field {
	return zap.ByteStrings(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Float32p(key string, val *float32) Field {
	return zap.Float32p(key, val)
}

func Float32s(key string, val []float32) Field {
	return zap.Float32s(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float64p(key string, val *float64) Field {
	return zap.Float64p(key, val)
}

func Float64s(key string, val []float64) Field {
	return zap.Float64s(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Intp(key string, val *int) Field {
	return zap.Intp(key, val)
}

func Ints(key string, val []int) Field {
	return zap.Ints(key, val)
}

func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func Int8p(key string, val *int8) Field {
	return zap.Int8p(key, val)
}

func Int8s(key string, val []int8) Field {
	return zap.Int8s(key, val)
}

func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int16p(key string, val *int16) Field {
	return zap.Int16p(key, val)
}

func Int16s(key string, val []int16) Field {
	return zap.Int16s(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int32p(key string, val *int32) Field {
	return zap.Int32p(key, val)
}

func Int32s(key string, val []int32) Field {
	return zap.Int32s(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int64p(key string, val *int64) Field {
	return zap.Int64p(key, val)
}

func Int64s(key string, val []int64) Field {
	return zap.Int64s(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uintp(key string, val *uint) Field {
	return zap.Uintp(key, val)
}

func Uints(key string, val []uint) Field {
	return zap.Uints(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Uint8p(key string, val *uint8) Field {
	return zap.Uint8p(key, val)
}

func Uint8s(key string, val []uint8) Field {
	return zap.Uint8s(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint16p(key string, val *uint16) Field {
	return zap.Uint16p(key, val)
}

func Uint16s(key string, val []uint16) Field {
	return zap.Uint16s(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint32p(key string, val *uint32) Field {
	return zap.Uint32p(key, val)
}

func Uint32s(key string, val []uint32) Field {
	return zap.Uint32s(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint64p(key string, val *uint64) Field {
	return zap.Uint64p(key, val)
}

func Uint64s(key string, val []uint64) Field {
	return zap.Uint64s(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func Stringp(key string, val *string) Field {
	return zap.Stringp(key, val)
}

func Strings(key string, val []string) Field {
	return zap.Strings(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) Field {
	return zap.Timep(key, val)
}

func Times(key string, val []time.Time) Field {
	return zap.Times(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) Field {
	return zap.Durationp(key, val)
}

func Durations(key string, val []time.Duration) Field {
	return zap.Durations(key, val)
}

// Stack adds the current stacktrace to the log
func Stack(key string) Field {
	return zap.Stack(key)
}
