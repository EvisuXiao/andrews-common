package utils

import (
	"reflect"
)

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func HasErr(err error) bool {
	return err != nil
}

func IsSimpleValue(value interface{}) bool {
	switch value.(type) {
	case bool, string, byte, rune,
		uint, uint16, uint32, uint64, uintptr,
		int, int8, int16, int64:
		return true
	}
	return false
}

func IsArrayOrMap(value interface{}) bool {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return true
	}
	return false
}
