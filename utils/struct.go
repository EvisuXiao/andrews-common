package utils

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StructToMap(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	tp := v.Type()
	if tp.Kind() != reflect.Struct {
		return m
	}
	num := tp.NumField()
	for i := 0; i < num; i++ {
		tpField := tp.Field(i)
		vField := v.Field(i)
		name := tpField.Name
		jsonTag := tpField.Tag.Get("json")
		jsonTagArr := strings.Split(jsonTag, ",")
		if !IsEmpty(jsonTagArr) && !IsEmpty(jsonTagArr[0]) {
			name = jsonTagArr[0]
		}
		if name == "-" || !vField.CanInterface() {
			continue
		}
		value := vField.Interface()
		if tpField.Anonymous {
			m = MergeMap(m, StructToMap(value))
			continue
		}
		_, isTime := value.(time.Time)
		subType := tpField.Type.Kind()
		if subType == reflect.Ptr {
			subType = tpField.Type.Elem().Kind()
		}
		if !isTime && subType == reflect.Struct {
			m[name] = StructToMap(vField)
			continue
		}
		m[name] = vField.Interface()
	}
	return m
}

func ArrayMapStringToStruct(m []map[string]string, data interface{}, tFormat ...string) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return
	}
	v = v.Elem()
	kind := v.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return
	}
	var item reflect.Value
	for _, mm := range m {
		subType := v.Type().Elem()
		if subType.Kind() == reflect.Ptr {
			item = reflect.New(subType.Elem())
		} else {
			item = reflect.New(subType)
		}
		MapStringToStruct(mm, item.Interface(), tFormat...)
		v.Set(reflect.Append(v, item))
	}
}

func MapStringToStruct(m map[string]string, data interface{}, tFormat ...string) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.CanAddr() {
		return
	}
	tp := v.Type()
	if tp.Kind() != reflect.Struct {
		return
	}
	datetimeFormat := DatetimeFormat
	dateFormat := DateFormat
	timeFormatLen := len(tFormat)
	if timeFormatLen > 0 {
		datetimeFormat = tFormat[0]
	}
	if timeFormatLen > 1 {
		dateFormat = tFormat[1]
	}
	num := tp.NumField()
	for i := 0; i < num; i++ {
		tpField := tp.Field(i)
		vField := v.Field(i)
		subKind := vField.Kind()
		if !vField.CanInterface() {
			continue
		}
		if tpField.Anonymous {
			if subKind == reflect.Ptr {
				if vField.IsNil() {
					vField.Set(reflect.New(vField.Type().Elem()))
				}
				MapStringToStruct(m, vField.Interface())
			} else if subKind == reflect.Struct && vField.CanAddr() {
				MapStringToStruct(m, vField.Addr().Interface())
			}
			continue
		}
		name := tpField.Name
		jsonTag := tpField.Tag.Get("json")
		jsonTagArr := strings.Split(jsonTag, ",")
		if !IsEmpty(jsonTagArr) && !IsEmpty(jsonTagArr[0]) {
			name = jsonTagArr[0]
		}
		if name == "-" {
			continue
		}
		if val, ok := m[name]; ok {
			if vField.CanSet() {
				if subKind == reflect.Ptr {
					vField = vField.Elem()
				}
				if _, ok := vField.Interface().(time.Time); ok {
					timeVal, err := time.ParseInLocation(datetimeFormat, val, time.Local)
					if HasErr(err) {
						timeVal, err = time.ParseInLocation(dateFormat, val, time.Local)
					}
					if !HasErr(err) {
						vField.Set(reflect.ValueOf(timeVal))
					}
					continue
				}
				setSimpleFieldValue(vField, val)
			}
		}
	}
}

func SetStructDefaultValue(s interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.CanAddr() {
		return
	}
	tp := v.Type()
	if tp.Kind() != reflect.Struct {
		return
	}
	num := tp.NumField()
	for i := 0; i < num; i++ {
		vField := v.Field(i)
		if vField.Kind() == reflect.Ptr {
			if vField.IsNil() {
				vField.Set(reflect.New(vField.Type().Elem()))
			}
			vField = vField.Elem()
		}
		if !vField.CanInterface() {
			return
		}
		if vField.Kind() == reflect.Struct {
			SetStructDefaultValue(vField.Addr().Interface())
			continue
		}
		if vField.CanSet() {
			defaultVal := tp.Field(i).Tag.Get("default")
			if IsEmpty(defaultVal) {
				continue
			}
			if vField.IsZero() {
				setSimpleFieldValue(vField, defaultVal)
			}
		}
	}
}

func setSimpleFieldValue(field reflect.Value, val string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		strVal, _ := strconv.ParseInt(val, 10, 64)
		field.SetInt(strVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		strVal, _ := strconv.ParseUint(val, 10, 64)
		field.SetUint(strVal)
	case reflect.Float32, reflect.Float64:
		strVal, _ := strconv.ParseFloat(val, 64)
		field.SetFloat(strVal)
	case reflect.Bool:
		strVal, _ := strconv.ParseBool(val)
		field.SetBool(strVal)
	}
}
