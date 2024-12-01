package model

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/now"

	"github.com/lookupearth/restful/field"
)

type FieldCase struct {
	ID       int64                  `default:"123" gorm:"primarykey"`
	Status   field.ExInt64          `default:"-578" gorm:"column:status2" operate:">="`
	Name     string                 `default:"555" gorm:"-"`
	NameTest string                 ``
	Time     *field.Time            `default:"2022-06-24 00:00:00" gorm:"-:migration"`
	Time2    *field.Time            `default:"now" gorm:"-:migration"`
	Time3    *field.Timestamp       `default:"now" gorm:"-:migration"`
	Test     map[string]interface{} ``
}

func TestField(t *testing.T) {
	structData := &FieldCase{}
	fmt.Println(structData)
	tm := time.Date(2022, 06, 24, 00, 00, 00, 00, time.Local)
	tm, _ = now.Parse("2022-06-24 00:00:00")
	cases := []struct {
		Field        string
		HaveDefault  bool
		DefaultValue interface{}
		ParseValue   string
		ParseResult  interface{}
		WhereValue   string
		Primarykey   bool
	}{
		{
			Field:        "ID",
			HaveDefault:  true,
			DefaultValue: int64(123),
			ParseValue:   "456",
			ParseResult:  int64(456),
			WhereValue:   "`id` = ?",
			Primarykey:   true,
		},
		{
			Field:        "Status",
			HaveDefault:  true,
			DefaultValue: field.ExInt64(-578),
			ParseValue:   "456",
			ParseResult:  field.ExInt64(456),
			WhereValue:   "`status2` >= ?",
			Primarykey:   false,
		},
		{
			Field:        "Name",
			HaveDefault:  true,
			DefaultValue: "555",
			ParseValue:   "777",
			ParseResult:  "777",
			WhereValue:   "`name` = ?",
			Primarykey:   false,
		},
		{
			Field:        "NameTest",
			HaveDefault:  false,
			DefaultValue: nil,
			ParseValue:   "888",
			ParseResult:  "888",
			WhereValue:   "`name_test` = ?",
			Primarykey:   false,
		},
		{
			Field:        "Time",
			HaveDefault:  true,
			DefaultValue: field.Time(tm),
			ParseValue:   "2022-06-24 00:00:00",
			ParseResult:  field.Time(tm),
			WhereValue:   "`time` = ?",
			Primarykey:   false,
		},
		{
			Field:        "Time2",
			HaveDefault:  true,
			DefaultValue: nil, // ms以下时间戳不可能对齐
			ParseValue:   "2022-06-24 00:00:00",
			ParseResult:  field.Time(tm),
			WhereValue:   "`time2` = ?",
			Primarykey:   false,
		},
		{
			Field:        "Time3",
			HaveDefault:  true,
			DefaultValue: nil, // ms以下时间戳不可能对齐
			ParseValue:   "1656000000",
			ParseResult:  field.Timestamp(tm),
			WhereValue:   "`time3` = ?",
			Primarykey:   false,
		},
		{
			Field:        "Test",
			HaveDefault:  false,
			DefaultValue: nil,
			ParseValue:   `{"a":1}`,
			ParseResult:  map[string]interface{}{"a": float64(1)},
			WhereValue:   "`test` = ?",
			Primarykey:   false,
		},
	}

	ctx := context.Background()
	st := reflect.TypeOf(structData).Elem()
	for _, c := range cases {
		fs, _ := st.FieldByName(c.Field)
		f := NewField(fs)
		if f.HaveDefaultValue() != c.HaveDefault {
			t.Errorf("field %s have default parse fail, expect: %t, got: %t", c.Field, c.HaveDefault, f.HaveDefaultValue())
		}
		v := f.GetDefaultValue(ctx)
		if v != c.DefaultValue && c.DefaultValue != nil {
			t.Errorf("field %s default parse fail, expect: %s, got: %s", c.Field, c.DefaultValue, v)
		}
		if len(c.ParseValue) > 0 {
			pv, err := f.Parse(c.ParseValue)
			if err != nil {
				t.Errorf("field %s parse fail, value=%s", c.Field, c.ParseValue)
			}
			if !reflect.DeepEqual(pv, c.ParseResult) {
				t.Errorf("field %s parse fail, expect: %s, got: %s", c.Field, c.ParseResult, pv)
			}
		}
		if f.Where() != c.WhereValue {
			t.Errorf("field %s where fail, expect: %s, got: %s", c.Field, c.WhereValue, f.Where())
		}
		if f.PrimaryKey != c.Primarykey {
			t.Errorf("field %s Primarykey fail, expect: %t, got: %t", c.Field, c.Primarykey, f.PrimaryKey)
		}
	}
}
