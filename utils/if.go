package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// callFn if args[i] == func, run it
func callFn(f interface{}) interface{} {
	if f != nil {
		t := reflect.TypeOf(f)
		if t.Kind() == reflect.Func && t.NumIn() == 0 {
			function := reflect.ValueOf(f)
			in := make([]reflect.Value, 0)
			out := function.Call(in)
			if num := len(out); num > 0 {
				list := make([]interface{}, num)
				for i, value := range out {
					list[i] = value.Interface()
				}
				if num == 1 {
					return list[0]
				}
				return list
			}
			return nil
		}
	}
	return f
}

// TestFnTime run func use time
func TestFnTime(f interface{}) string {
	start := time.Now()
	callFn(f)
	end := time.Now()
	vf := reflect.ValueOf(f)
	str := fmt.Sprintf("[%s] runtime: %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
	fmt.Println(str)
	return str
}

// If - (a ? b : c) Or (a && b)
func If(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	var trueVal = args[1]
	var falseVal interface{}
	if len(args) > 2 {
		falseVal = args[2]
	} else {
		falseVal = nil
	}
	if condition == nil {
		return callFn(falseVal)
	} else if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(falseVal)
		}
	} else if IsEmpty(condition) {
		return callFn(falseVal)
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return callFn(trueVal)
}

// Or - (a || b)
func Or(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	if condition == nil {
		return callFn(args[1])
	}
	if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(args[1])
		}
	} else if IsEmpty(condition) {
		return callFn(args[1])
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return condition
}
