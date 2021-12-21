package utils

import (
	"fmt"
	"reflect"
)

func SliceIntJoin(a []int, sep string) string {
	length := len(a)
	if IsEmpty(length) {
		return ""
	}
	var s string
	for i := 0; i < length; i++ {
		s += fmt.Sprintf("%s%d", sep, a[i])
	}
	return s[1:]
}

func InSlice(needle interface{}, haystack interface{}) bool {
	return SearchSliceIndex(needle, haystack) > -1
}

func SearchSliceIndex(needle interface{}, haystack interface{}) int {
	haystackRv := reflect.ValueOf(haystack)
	switch haystackRv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < haystackRv.Len(); i++ {
			if reflect.DeepEqual(needle, haystackRv.Index(i).Interface()) == true {
				return i
			}
		}
	default:
		panic("Func SearchSliceIndex: parameter haystack should be array or slice")
	}
	return -1
}

type SliceIntCompare struct {
	Add   []int
	Sub   []int
	Assoc []int
	Merge []int
}

func CompareSliceInt(before, after []int) *SliceIntCompare {
	compare := &SliceIntCompare{}
	for _, v := range after {
		SliceAddIntItem(&compare.Merge, v)
		if InSlice(v, before) {
			SliceAddIntItem(&compare.Assoc, v)
		} else {
			SliceAddIntItem(&compare.Add, v)
		}
	}
	for _, v := range before {
		SliceAddIntItem(&compare.Merge, v)
		if !InSlice(v, after) {
			SliceAddIntItem(&compare.Sub, v)
		}
	}
	return compare
}

func SliceAddIntItem(slice *[]int, item int) {
	if !InSlice(item, *slice) {
		*slice = append(*slice, item)
	}
}

func SliceAddStringItem(slice *[]string, item string) {
	if !InSlice(item, *slice) {
		*slice = append(*slice, item)
	}
}

func SliceDeleteIntItem(slice []int, idx int) []int {
	return append(slice[:idx], slice[idx+1:]...)
}

func SliceDeleteStringItem(slice []string, idx int) []string {
	return append(slice[:idx], slice[idx+1:]...)
}

func SliceIntToInt64(val []int) []int64 {
	var m []int64
	for _, v := range val {
		m = append(m, int64(v))
	}
	return m
}

func SliceInt64ToInt(val []int64) []int {
	var m []int
	for _, v := range val {
		m = append(m, int(v))
	}
	return m
}

func SliceInterfaceToInt(val []interface{}) []int {
	var m []int
	for _, v := range val {
		m = append(m, int(v.(float64)))
	}
	return m
}

func SliceInterfaceToString(val []interface{}) []string {
	var m []string
	for _, v := range val {
		m = append(m, v.(string))
	}
	return m
}
