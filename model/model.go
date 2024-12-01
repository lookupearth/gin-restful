package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Model model解析类
type Model struct {
	ModelInterface interface{}
	ModelType      reflect.Type
	// DB 主键标识
	PrimaryKey string

	// struct名到json字段
	Name2Json map[string]string
	// struct名到db列
	Name2Column map[string]string
	// json字段到struct名
	Json2Name map[string]string
	// db列到struct名
	Column2Name map[string]string
	// struct名到db列
	Name2Field map[string]*Field
}

// NewModel NewParser Model 实例化，确保尽在启动阶段调用，而不会在请求处理阶段调用
func NewModel(model interface{}) *Model {
	mt := reflect.TypeOf(model)
	mv := reflect.ValueOf(model)
	if mt.Kind() == reflect.Ptr {
		mt = mt.Elem()
	}
	if mv.Kind() == reflect.Ptr {
		mv = mv.Elem()
	}
	p := &Model{
		ModelInterface: model,
		ModelType:      mt,
	}
	p.Init()
	return p
}

func (model *Model) Init() {
	model.Name2Json = make(map[string]string)
	model.Name2Column = make(map[string]string)
	model.Json2Name = make(map[string]string)
	model.Column2Name = make(map[string]string)
	model.Name2Field = make(map[string]*Field)

	t := model.ModelType
	// 遍历结构体中所有字段
	for i := 0; i < t.NumField(); i++ {
		// 单个结构体字段
		f := t.Field(i)
		name := f.Name
		field := NewField(f)
		model.Name2Field[name] = field
		if len(field.Json.Name) > 0 {
			model.Name2Json[name] = field.Json.Name
			model.Json2Name[field.Json.Name] = name
		}
		if len(field.Gorm.Column) > 0 {
			model.Name2Column[name] = field.Gorm.Column
			model.Column2Name[field.Gorm.Column] = name
		}
		if field.PrimaryKey && len(model.PrimaryKey) == 0 {
			model.PrimaryKey = field.Gorm.Column
		} else if field.PrimaryKey && len(model.PrimaryKey) != 0 {
			panic(fmt.Sprintf("model <%s> can only have one PrimaryKey", model.ModelType.Name()))
		}
	}
}

// New 实例化具体Model
func (model *Model) New() interface{} {
	return reflect.New(model.ModelType).Interface()
}

// NewSlice 实例化具体Model的切片，用于列表查询
func (model *Model) NewSlice() interface{} {
	slice := reflect.MakeSlice(reflect.SliceOf(model.ModelType), 0, 0)
	slicePtr := reflect.New(slice.Type())
	slicePtr.Elem().Set(slice)
	return slicePtr.Interface()
}

// Parse 注意返回的是model的指针
func (model *Model) Parse(b []byte) (interface{}, error) {
	data := model.New()
	err := json.Unmarshal(b, data)
	return data, err
}

// ParseFromQuery 注意返回的是model的指针
func (model *Model) ParseFromQuery(query map[string]string) (interface{}, error) {
	data := model.New()
	dv := reflect.ValueOf(data).Elem()
	for k, v := range query {
		if name, ok := model.Json2Name[k]; ok {
			value := v
			vv, err := model.Name2Field[name].Parse(value)
			if err != nil {
				return nil, err
			}
			fv := dv.FieldByName(name)
			makePtr(fv).Set(reflect.ValueOf(vv))
		}
	}
	return data, nil
}

// ParseDefault 解析默认值，添加的默认值会被记录到input中
func (model *Model) ParseDefault(ctx context.Context, data interface{}, input map[string]interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(data))
	for name, field := range model.Name2Field {
		fv := value.FieldByName(name)
		if jsonKey, ok := model.Name2Json[name]; ok {
			if _, ok2 := input[jsonKey]; !ok2 {
				if field.HaveDefaultValue() {
					makePtr(fv).Set(reflect.ValueOf(field.GetDefaultValue(ctx)))
					input[jsonKey] = 1
				}
			}
		} else {
			if field.HaveDefaultValue() {
				makePtr(fv).Set(reflect.ValueOf(field.GetDefaultValue(ctx)))
			}
		}
	}
	return nil
}

// ParseDefaultWithKeys 按固定输入key解析默认值，添加的默认值会被记录到input中
func (model *Model) ParseDefaultWithKeys(ctx context.Context, data interface{}, keys []string, input map[string]interface{}) error {
	if keys == nil {
		return nil
	}
	value := reflect.Indirect(reflect.ValueOf(data))
	for _, name := range keys {
		field, ok := model.Name2Field[name]
		if !ok {
			continue
		}
		fv := value.FieldByName(name)
		if jsonKey, ok := model.Name2Json[name]; ok {
			if _, ok2 := input[jsonKey]; !ok2 {
				if field.HaveDefaultValue() {
					makePtr(fv).Set(reflect.ValueOf(field.GetDefaultValue(ctx)))
					input[jsonKey] = 1
				}
			}
		} else {
			if field.HaveDefaultValue() {
				makePtr(fv).Set(reflect.ValueOf(field.GetDefaultValue(ctx)))
			}
		}
	}
	return nil
}

func (model *Model) CheckPrimaryKey() error {
	if len(model.PrimaryKey) == 0 {
		return errors.New("model need a PrimaryKey")
	}
	return nil
}

// ParsePrimaryKey 注意返回的是model的指针
func (model *Model) ParsePrimaryKey(primaryKey string) (interface{}, error) {
	if err := model.CheckPrimaryKey(); err != nil {
		return nil, err
	}
	if name, ok := model.Column2Name[model.PrimaryKey]; ok {
		field, err := model.Name2Field[name].Parse(primaryKey)
		if err != nil {
			return nil, err
		}
		return field, nil
	}
	return nil, errors.New("primaryKey need a gorm column")
}

func (model *Model) FieldNames(data map[string]interface{}) []string {
	names := make([]string, 0)
	for k := range data {
		if name, ok := model.Json2Name[k]; ok {
			names = append(names, name)
		}
	}
	return names
}

// Where 获取字段的where条件，key为 db 中的 列名
func (model *Model) Where(query *gorm.DB, key string, value interface{}) *gorm.DB {
	if name, ok := model.Json2Name[key]; ok {
		field := model.Name2Field[name]
		query = query.Where(field.Where(), field.WhereValue(value))
	}
	return query
}
