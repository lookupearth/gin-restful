package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Unmarshaler interface {
	UnmarshalString(string) error
}

// parseValue 将字符串v解析为t类型的指针，返回nil, nil，说明不是基础类型
func parseValue(t reflect.Type, v string) (interface{}, error) {
	switch t.Kind() {
	case reflect.Bool:
		rv, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s for bool, got error: %v", v, err)
		}
		return reflect.ValueOf(rv).Convert(t).Interface(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s for int, got error: %v", v, err)
		}
		return reflect.ValueOf(rv).Convert(t).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s for uint, got error: %v", v, err)
		}
		return reflect.ValueOf(rv).Convert(t).Interface(), nil
	case reflect.Float32, reflect.Float64:
		rv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s for float, got error: %v", v, err)
		}
		return reflect.ValueOf(rv).Convert(t).Interface(), nil
	case reflect.String:
		return reflect.ValueOf(v).Convert(t).Interface(), nil
	default:
		return nil, nil
	}
}

// parseJson 将字符串v解析为t类型的指针，返回nil, nil，说明不是map/slice类型
func parseJson(t reflect.Type, v string) (interface{}, error) {
	switch t.Kind() {
	case reflect.Map, reflect.Slice:
		fieldPtr := reflect.New(t)
		if err := json.Unmarshal([]byte(v), fieldPtr.Interface()); err != nil {
			return nil, fmt.Errorf("failed to parse %s for json, got error: %v", v, err)
		}
		return fieldPtr.Elem().Interface(), nil
	case reflect.Struct:
		fieldPtr := reflect.New(t)
		if err := json.Unmarshal([]byte(v), fieldPtr.Interface()); err != nil {
			return nil, nil
		}
		return fieldPtr.Elem().Interface(), nil
	default:
		return nil, nil
	}
}

func isTime(t reflect.Type) bool {
	timeType := false
	if t.Kind() == reflect.Struct {
		if _, ok := reflect.New(t).Elem().Interface().(time.Time); ok {
			timeType = true
		} else if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
			timeType = true
		}
	}
	return timeType
}

func makePtr(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		if value.Kind() == reflect.Ptr && value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}
	return value
}
