package utils

import (
	"reflect"
)

func MergeMap(map1, map2 map[string]interface{}) map[string]interface{} {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func GetMapValueString(m, k interface{}) string {
	rv := reflect.ValueOf(m)
	if rv.Kind() != reflect.Map {
		return ""
	}
	v := rv.MapIndex(reflect.ValueOf(k))
	if v.IsValid() && v.Kind() == reflect.String {
		return v.String()
	}
	return ""
}

func InMap(needle interface{}, haystack interface{}) bool {
	return !IsEmpty(SearchMapKey(needle, haystack))
}

func SearchMapKey(needle interface{}, haystack interface{}) interface{} {
	haystackRv := reflect.ValueOf(haystack)
	switch haystackRv.Kind() {
	case reflect.Map:
		for i := 0; i < haystackRv.Len(); i++ {
			if reflect.DeepEqual(needle, haystackRv.Index(i).Interface()) == true {
				return haystackRv.Index(i).Interface()
			}
		}
	default:
		panic("Func SearchMapKey: parameter haystack should be map")
	}
	return nil
}
