package model

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/now"
)

type IFieldDefault interface {
	GetDefault(context.Context, interface{}) interface{}
}

type IFieldDefaultPrepare interface {
	GetDefaultPrepare(reflect.Type, string) (interface{}, error)
}

type Default struct {
	Value          string
	ValueInterface interface{}
	DefaultFunc    DefaultFunc

	haveDefault bool
	err         error
	fieldType   reflect.Type
}

func NewDefault(field reflect.StructField) *Default {
	fieldType := field.Type
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	tag, ok := field.Tag.Lookup("default")
	cols := strings.Split(tag, ":")
	funcName := cols[0]
	params := ""
	if len(cols) > 1 {
		params = strings.Join(cols[1:], ":")
	}
	f := &Default{
		Value:       tag,
		haveDefault: ok,
		fieldType:   fieldType,
	}
	if f.haveDefault {
		fieldValue := reflect.New(fieldType).Elem().Interface()
		// 优先尝试解析一下
		if f.haveDefault {
			f.parseDefault(fieldType, tag)
		}
		// 不论解析成功与否，尝试验证是否实现 IFieldDefault 或 注册有相关函数
		if defaultFunc, ok := fieldValue.(IFieldDefault); ok {
			if prepareFunc, ok := fieldValue.(IFieldDefaultPrepare); ok {
				f.processDefaultFunc(tag, prepareFunc.GetDefaultPrepare, defaultFunc.GetDefault)
			} else {
				f.processDefaultFunc(tag, nil, defaultFunc.GetDefault)
			}
		}
		if len(funcName) > 0 && GetDefaultFunc(fieldValue, funcName) != nil {
			defaultFunc := GetDefaultFunc(fieldValue, funcName)
			prepareFunc := GetDefaultPrepareFunc(fieldValue, funcName)
			f.processDefaultFunc(params, prepareFunc, defaultFunc)
		}

		if f.err != nil {
			panic(f.err)
		}
	}
	return f
}

func (d *Default) processDefaultFunc(params string, dpf DefaultPrepareFunc, df DefaultFunc) {
	if df == nil && dpf == nil {
		return
	}
	var err error = nil
	var prepareValue interface{}
	if dpf != nil {
		prepareValue, err = dpf(d.fieldType, params)
	} else if d.ValueInterface == nil {
		prepareValue = params
	}
	// 若返回的不是ignoreErr，触发panic，否则相当于DefaultFunc无效
	if err != nil {
		if err.Error() == "ignore" {
			if d.ValueInterface == nil && d.DefaultFunc == nil {
				d.haveDefault = false
			}
		} else {
			d.err = err
		}
	} else {
		// 返回的是nil，说明DefaultFunc生效了，不需要err了
		d.ValueInterface = prepareValue
		d.DefaultFunc = df
		d.err = nil
	}
}

func (d *Default) HaveDefault() bool {
	return d.haveDefault
}

// GetValue 获取默认值，优先使用 DefaultFunc
func (d *Default) GetValue(ctx context.Context) interface{} {
	if d.DefaultFunc != nil {
		return d.DefaultFunc(ctx, d.ValueInterface)
	}
	return d.ValueInterface
}

func (d *Default) parseDefault(fieldType reflect.Type, value string) {
	value = strings.TrimSpace(value)
	// 解析特殊类型
	fieldPtr := reflect.New(fieldType)
	unmarshaler, isUnmarshaler := fieldPtr.Interface().(Unmarshaler)
	if isUnmarshaler {
		// 由于无法保证不是引用对象，因此必须通过函数方式生成
		err := unmarshaler.UnmarshalString(value)
		if err == nil {
			d.ValueInterface = value
			d.DefaultFunc = func(ctx context.Context, v interface{}) interface{} {
				vv := reflect.New(fieldType)
				ju, _ := vv.Interface().(Unmarshaler)
				_ = ju.UnmarshalString(v.(string))
				return vv.Elem().Interface()
			}
			return
		}
	} else if isTime(fieldType) {
		_, err := now.Parse(value)
		if err != nil {
			d.err = fmt.Errorf("failed to parse %s for time", d.Value)
		}
		d.ValueInterface = value
		d.DefaultFunc = func(ctx context.Context, v interface{}) interface{} {
			tmp, _ := now.Parse(v.(string))
			vv := reflect.New(fieldType).Elem()
			vv.Set(reflect.ValueOf(tmp).Convert(fieldType))
			return vv.Interface()
		}
		return
	} else {
		// 确认是否可以当做json处理，由于map/slice等为引用对象，需要每次创建一个新的
		err := json.Unmarshal([]byte(value), fieldPtr.Interface())
		if err == nil {
			d.ValueInterface = []byte(value)
			d.DefaultFunc = func(ctx context.Context, v interface{}) interface{} {
				vv := reflect.New(fieldType)
				_ = json.Unmarshal(v.([]byte), vv.Interface())
				return vv.Elem().Interface()
			}
			return
		}
	}
	// 解析基础类型
	v, err := parseValue(fieldType, value)
	if err != nil {
		d.err = err
		return
	}
	if v != nil {
		d.ValueInterface = v
		return
	}
	d.err = fmt.Errorf("failed to parse %s for %s", d.Value, fieldType)
}
