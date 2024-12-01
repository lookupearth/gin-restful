// Package model
package model

import (
	"reflect"
	"testing"
)

type JsonCase struct {
	ID     int64  `gorm:"column:id;primaryKey;->" json:"id"`            // 自增id
	Status int64  `gorm:"column:status;default:2" json:"-" operate:"="` // 状态
	Name   string `gorm:"column:name;default:app" json:"" operate:"="`  // 名称
	Name2  string `gorm:"column:name;default:app" operate:"="`          // 名称
}

func TestJson(t *testing.T) {
	structData := &JsonCase{}
	cases := []struct {
		Field string
		Key   string
	}{
		{
			Field: "ID",
			Key:   "id",
		},
		{
			Field: "Status",
			Key:   "",
		},
		{
			Field: "Name",
			Key:   "Name",
		},
		{
			Field: "Name2",
			Key:   "Name2",
		},
	}
	st := reflect.TypeOf(structData).Elem()
	for _, c := range cases {
		f, _ := st.FieldByName(c.Field)
		j := NewJson(f)
		if j.Name != c.Key {
			t.Errorf("field %s json key parse fail", c.Field)
		}
	}

}
