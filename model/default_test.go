package model

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/lookupearth/restful/field"
)

type DefaultTestStruct struct {
	Name string
}

func (d DefaultTestStruct) GetDefault(ctx context.Context, v interface{}) interface{} {
	return DefaultTestStruct{
		Name: v.(string),
	}
}

func (d DefaultTestStruct) GetDefaultPrepare(t reflect.Type, params string) (interface{}, error) {
	if params == "123" {
		return nil, errors.New("123")
	}
	return params, nil
}

type DefaultCase struct {
	ID      int64             `default:"123"`
	Status  *field.ExInt64    `default:"-578"`
	Name    string            `default:"555" json:"name"`
	Name2   string            ``
	Name3   *field.ExString   `default:""`
	ID2     float64           `default:"12.3"`
	ID3     uint32            `default:"123" json:"-"`
	Status2 bool              `default:"1"`
	Status3 bool              `default:"false"`
	Time    time.Time         `default:"2022-06-24 00:00:00"`
	Time2   *field.Time       `default:"2022-06-24 00:00:00"`
	Test    DefaultTestStruct `default:"1234"`
	Test2   *string           `default:"test:1234"`
	Test3   string            `default:"test2:777"`
	ID4     int64             `default:"test:123456"`
	Obj     field.JSONObject  `default:"{\"a\":1}"`
}

func TestDefault(t *testing.T) {
	RegisterDefaultPrepareFunc(func(t reflect.Type, params string) (interface{}, error) {
		if params == "test:123" {
			return nil, errors.New("123")
		}
		return params, nil
	}, "test", "")
	RegisterDefaultFunc(func(ctx context.Context, v interface{}) interface{} {
		return v
	}, "test", "")
	RegisterDefaultPrepareFunc(func(t reflect.Type, params string) (interface{}, error) {
		fmt.Println(params)
		if params == "test2:777" {
			return nil, &field.IgnoreError{}
		}
		return params, nil
	}, "test2", "")
	RegisterDefaultFunc(func(ctx context.Context, v interface{}) interface{} {
		return int64(123)
	}, "test", int64(1))
	structData := &DefaultCase{}
	fmt.Println(structData)
	tm := time.Date(2022, 06, 24, 00, 00, 00, 00, time.Local)
	cases := []struct {
		Field       string
		HaveDefault bool
		Value       interface{}
	}{
		{
			Field:       "ID",
			HaveDefault: true,
			Value:       int64(123),
		},
		{
			Field:       "Status",
			HaveDefault: true,
			Value:       field.ExInt64(-578),
		},
		{
			Field:       "Name",
			HaveDefault: true,
			Value:       "555",
		},
		{
			Field:       "Name2",
			HaveDefault: false,
			Value:       nil,
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
		{
			Field:       "Obj",
			HaveDefault: true,
			Value:       field.JSONObject(map[string]interface{}{"a": float64(1)}),
		},
	}
	ctx := context.Background()
	st := reflect.TypeOf(structData).Elem()
	sv := reflect.ValueOf(structData).Elem()
	for _, c := range cases {
		f, _ := st.FieldByName(c.Field)
		d := NewDefault(f)
		if d.HaveDefault() != c.HaveDefault {
			t.Errorf("field %s have default parse fail, expect: %t, got: %t", c.Field, c.HaveDefault, d.HaveDefault())
		}
		v := d.GetValue(ctx)
		if !reflect.DeepEqual(v, c.Value) {
			t.Errorf("field %s default parse fail, expect: %s, got: %s", c.Field, c.Value, v)
		}
		if d.HaveDefault() {
			fv := sv.FieldByName(c.Field)
			fmt.Println(c.Field, fv.CanAddr(), fv.Kind(), v, reflect.ValueOf(v))
			for fv.Kind() == reflect.Ptr {
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					fv.Set(reflect.New(fv.Type().Elem()))
				}
				fv = fv.Elem()
			}
			fv.Set(reflect.ValueOf(v))
		}
	}
	fmt.Println(structData)
	fmt.Println(structData.Name3)
}
