package model

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jinzhu/now"
)

type Field struct {
	FieldType reflect.Type

	PrimaryKey bool

	JsonKey string // 空表示不能从json读写
	DBKey   string // 空表示不与数据库交互

	Gorm    *Gorm
	Json    *Json
	Default *Default
	Operate *Operate
}

func NewField(field reflect.StructField) *Field {
	fieldType := field.Type
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	instance := &Field{
		FieldType:  fieldType,
		PrimaryKey: false,
		Json:       NewJson(field),
		Gorm:       NewGorm(field),
		Default:    NewDefault(field),
		Operate:    NewOperate(field),
	}
	instance.JsonKey = instance.Json.Name
	instance.DBKey = instance.Gorm.Column

	// 主键
	if _, ok := instance.Gorm.Tags["PRIMARYKEY"]; ok {
		instance.PrimaryKey = true
	}

	return instance
}

// New 实例化具体Field
func (field *Field) New() interface{} {
	return reflect.New(field.FieldType).Interface()
}

// Parse 解析一个field，返回对应类型的数据
func (field *Field) Parse(value string) (interface{}, error) {
	// 解析特殊类型
	fieldPtr := reflect.New(field.FieldType)
	unmarshaler, isUnmarshaler := fieldPtr.Interface().(Unmarshaler)
	if isUnmarshaler {
		err := unmarshaler.UnmarshalString(value)
		if err != nil {
			return nil, field.parseError(fmt.Errorf("failed to parse %s for %s", value, field.FieldType))
		}
		return fieldPtr.Elem().Interface(), nil
	} else if isTime(field.FieldType) {
		tmp, err := now.Parse(value)
		if err != nil {
			return nil, field.parseError(fmt.Errorf("failed to parse %s for time", value))
		}
		vv := reflect.New(field.FieldType).Elem()
		vv.Set(reflect.ValueOf(tmp).Convert(field.FieldType))
		return vv.Interface(), nil
	}
	// 解析基础类型
	v, err := parseValue(field.FieldType, value)
	if err != nil {
		return nil, field.parseError(err)
	}
	if v != nil {
		return v, nil
	}
	// 解析json类型
	v, err = parseJson(field.FieldType, value)
	if err != nil {
		return nil, field.parseError(err)
	}
	if v != nil {
		return v, nil
	}
	return nil, field.parseError(fmt.Errorf("failed to parse %s for %s", value, field.FieldType))
}

func (field Field) parseError(err error) error {
	return fmt.Errorf("%s parse failed, %s", field.Json.Name, err.Error())
}

func (field Field) Where() string {
	return fmt.Sprintf("`%s` %s ?", field.Gorm.Column, field.Operate.Operate)
}

func (field Field) WhereValue(value interface{}) interface{} {
	return field.Operate.Value(value)
}

func (field *Field) HaveDefaultValue() bool {
	return field.Default.HaveDefault()
}

func (field *Field) GetDefaultValue(ctx context.Context) interface{} {
	return field.Default.GetValue(ctx)
}

func (field Field) ReadOnly() bool {
	return field.Json.ReadOnly
}
