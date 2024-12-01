package model

import (
	"context"
	"reflect"
)

type DefaultPrepareFunc func(reflect.Type, string) (interface{}, error)
type DefaultFunc func(context.Context, interface{}) interface{}

var defaultPrepareFuncs map[reflect.Type]map[string]DefaultPrepareFunc
var defaultFuncs map[reflect.Type]map[string]DefaultFunc

// RegisterDefaultPrepareFunc 注册默认值预处理函数
func RegisterDefaultPrepareFunc(fn DefaultPrepareFunc, tag string, fields ...interface{}) {
	if defaultPrepareFuncs == nil {
		defaultPrepareFuncs = make(map[reflect.Type]map[string]DefaultPrepareFunc)
	}
	for _, field := range fields {
		fieldType := reflect.TypeOf(field)
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if _, ok := defaultFuncs[fieldType]; !ok {
			defaultPrepareFuncs[fieldType] = make(map[string]DefaultPrepareFunc)
		}
		defaultPrepareFuncs[fieldType][tag] = fn
	}
}

// RegisterDefaultFunc 注册默认值处理函数
func RegisterDefaultFunc(fn DefaultFunc, tag string, fields ...interface{}) {
	if defaultFuncs == nil {
		defaultFuncs = make(map[reflect.Type]map[string]DefaultFunc)
	}
	for _, field := range fields {
		fieldType := reflect.TypeOf(field)
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if _, ok := defaultFuncs[fieldType]; !ok {
			defaultFuncs[fieldType] = make(map[string]DefaultFunc)
		}
		defaultFuncs[fieldType][tag] = fn
	}
}

// GetDefaultPrepareFunc 获取默认值预处理函数
func GetDefaultPrepareFunc(field interface{}, tag string) DefaultPrepareFunc {
	if defaultPrepareFuncs == nil {
		return nil
	}
	fieldType := reflect.TypeOf(field)
	if tagMap, ok := defaultPrepareFuncs[fieldType]; ok {
		if fn, ok := tagMap[tag]; ok {
			return fn
		}
	}
	return nil
}

// GetDefaultFunc 获取默认值处理函数
func GetDefaultFunc(field interface{}, tag string) DefaultFunc {
	if defaultFuncs == nil {
		return nil
	}
	fieldType := reflect.TypeOf(field)
	if tagMap, ok := defaultFuncs[fieldType]; ok {
		if fn, ok := tagMap[tag]; ok {
			return fn
		}
	}
	return nil
}
