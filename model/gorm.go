package model

import (
	"reflect"

	"gorm.io/gorm/schema"
)

type Gorm struct {
	Column string
	Tags   map[string]string
}

// NewGorm 创建一个gorm tag分析的struct
// 注意这里只使用了默认的列名策略，支持自定义策略，需要透传db对象，过于复杂
func NewGorm(field reflect.StructField) *Gorm {
	tag := field.Tag.Get("gorm")
	tagMap := schema.ParseTagSetting(tag, ";")
	name := ""
	if column, ok := tagMap["COLUMN"]; ok {
		if column == "-" {
			column = ""
		}
		name = column
	} else {
		ns := schema.NamingStrategy{}
		name = ns.ColumnName("", field.Name)
	}

	f := &Gorm{
		Column: name,
		Tags:   tagMap,
	}
	return f
}
