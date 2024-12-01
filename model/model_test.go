package model

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/lookupearth/restful/field"
)

type Activity struct {
	ID      int64           `gorm:"column:id;primaryKey;->" json:"id"`                 // 自增id
	Status  *field.ExInt64  `gorm:"column:status;default:2" json:"status" default:"1"` // 状态
	Status2 field.ExInt64   `gorm:"column:status2;default:2" json:"status2"`           // 状态
	Name    string          `gorm:"column:name;default:app" json:"name"`               // 名称
	Name2   *field.ExString `default:""`
}

func TestModelParse(t *testing.T) {
	model := NewModel(&Activity{})
	_ = model.New()
	b := []byte(`{"name": "test"}`)
	data, err := model.Parse(b)
	if err != nil {
		t.Errorf("model parse fail, error=%v", err)
	}
	value := data.(*Activity)
	if value.Name != "test" {
		t.Errorf("Activity Name error")
	}
	if value.Status2 != 0 {
		t.Errorf("Activity Status2 error")
	}
}

func TestModelParseFromQuery(t *testing.T) {
	model := NewModel(&Activity{})
	_ = model.New()
	query := make(map[string]string)
	query["status"] = "123"
	query["status2"] = "456"
	data, err := model.ParseFromQuery(query)
	if err != nil {
		t.Errorf("model ParseFromQuery fail, error=%v", err)
	}
	if *data.(*Activity).Status != 123 {
		t.Errorf("model status value fail")
	}
	if data.(*Activity).Status2 != 456 {
		t.Errorf("model status2 value fail")
	}
}

func TestModelParseDefault(t *testing.T) {
	ctx := context.Background()
	model := NewModel(&DefaultCase{})
	b := []byte(`{"ID": 1, "ID3": 4, "name": "test"}`)
	data, err := model.Parse(b)
	if err != nil {
		t.Errorf("model parse fail, error=%v", err)
	}
	var rawData map[string]interface{}
	_ = json.Unmarshal(b, &rawData)
	if err := model.ParseDefault(ctx, data, rawData); err != nil {
		t.Errorf("model ParseDefault fail, error=%v", err)
	}
	structData := data.(*DefaultCase)
	tm := time.Date(2022, 06, 24, 00, 00, 00, 00, time.Local)
	cases := []struct {
		Field       string
		HaveDefault bool
		Value       interface{}
	}{
		{
			Field:       "ID",
			HaveDefault: true,
			Value:       int64(1),
		},
		{
			Field:       "Status",
			HaveDefault: true,
			Value:       field.ExInt64(-578),
		},
		{
			Field:       "Name",
			HaveDefault: true,
			Value:       "test",
		},
		{
			Field:       "Name2",
			HaveDefault: false,
			Value:       "",
		},
		{
			Field:       "Name3",
			HaveDefault: true,
			Value:       field.ExString(""),
		},
		{
			Field:       "ID2",
			HaveDefault: true,
			Value:       12.3,
		},
		{
			Field:       "ID3",
			HaveDefault: true,
			Value:       uint32(123),
		},
		{
			Field:       "Status2",
			HaveDefault: true,
			Value:       true,
		},
		{
			Field:       "Status3",
			HaveDefault: true,
			Value:       false,
		},
		{
			Field:       "Time",
			HaveDefault: true,
			Value:       tm,
		},
		{
			Field:       "Time2",
			HaveDefault: true,
			Value:       field.Time(tm),
		},
		{
			Field:       "Test",
			HaveDefault: true,
			Value: DefaultTestStruct{
				Name: "1234",
			},
		},
		{
			Field:       "Test2",
			HaveDefault: true,
			Value:       "1234",
		},
		{
			Field:       "Test3",
			HaveDefault: true,
			Value:       "test2:777",
		},
		{
			Field:       "ID4",
			HaveDefault: true,
			Value:       int64(123),
		},
	}
	st := reflect.ValueOf(structData).Elem()
	for _, c := range cases {
		f := st.FieldByName(c.Field)
		for f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		if !reflect.DeepEqual(f.Interface(), c.Value) {
			t.Errorf("field %s default parse fail, expect: %s, got: %s", c.Field, f.Interface(), c.Value)
		}
	}
}
