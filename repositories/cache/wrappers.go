package cache

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/kgjoner/cornucopia/v2/utils/hash"
)

// Check for a cached result of F, if no hit, run it. F must return (*R, error).
func RunWithCache[R any, F any](q DAO, duration time.Duration, fn F) F {
	fnName := getFuncName(fn)

	wrapped := func(args ...any) (result R, err error) {
		key := fnName + ":" + hash.From(args...)

		err = q.GetJSON(key, &result)
		if err != ErrNil {
			return result, err
		}

		values := reflect.ValueOf(fn).Call(convertToReflectValues(args))
		resV, errV := values[0], values[1]
		if !errV.IsNil() {
			return result, errV.Interface().(error)
		}

		if !resV.IsNil() {
			result = reflect.Indirect(resV).Interface().(R)
			err = q.CacheJSON(key, result, duration)
		} else {
			err = q.CacheJSON(key, nil, duration)
		}

		return result, err
	}

	// Return a generic function of the same type as the input
	return reflect.MakeFunc(reflect.TypeOf(fn), func(in []reflect.Value) []reflect.Value {
		args := convertFromReflectValues(in)

		out, err := wrapped(args...)

		outV := reflect.ValueOf(out)
		if outV.IsZero() {
			outV = reflect.Zero(reflect.TypeOf(fn).Out(0))
		} else {
			outV = reflect.ValueOf(&out)
		}

		errV := reflect.ValueOf(err)
		if err == nil {
			errV = reflect.Zero(reflect.TypeOf(fn).Out(1))
		}

		return []reflect.Value{outV, errV}
	}).Interface().(F)
}

func getFuncName(fn any) string {
	pc := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	if pc == nil {
		return "unknown"
	}

	fullName := pc.Name()
	parts := strings.Split(fullName, "/")
	return parts[len(parts)-1]
}

// convertToReflectValues converts []any to []reflect.Value
func convertToReflectValues(args []any) []reflect.Value {
	values := make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}
	return values
}

// convertFromReflectValues converts []reflect.Value to []any
func convertFromReflectValues(values []reflect.Value) []any {
	args := make([]any, len(values))
	for i, val := range values {
		args[i] = val.Interface()
	}
	return args
}
