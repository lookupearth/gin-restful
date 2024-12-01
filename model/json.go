// Package model
package model

import (
	"reflect"
	"strings"
)

type Json struct {
	Name     string
	ReadOnly bool
}

func NewJson(field reflect.StructField) *Json {
	cols := strings.Split(field.Tag.Get("json"), ",")
	jsonKey := cols[0]
	if jsonKey == "" {
		jsonKey = field.Name
	} else if jsonKey == "-" {
		jsonKey = ""
	}
	readOnly := false
	for _, col := range cols {
		if strings.ToLower(col) == "readonly" {
			readOnly = true
			break
		}
	}
	f := &Json{
		Name:     jsonKey,
		ReadOnly: readOnly,
	}
	return f
}
