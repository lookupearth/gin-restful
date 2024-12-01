package field

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"time"
)

type Time time.Time

// Scan implements the Scanner interface.
// 将原始数据进行Struct结构改造
func (t *Time) Scan(value interface{}) error {
	// set default value
	*t = Time(time.Time{})

	if value == nil {
		return nil
	}

	vv, ok := value.(time.Time)
	if !ok {
		vt := reflect.TypeOf(value)
		if vt.Kind() == reflect.Ptr {
			vt = vt.Elem()
		}
		return errors.New("type error, expect time.Time received " + vt.String())
	}

	*t = Time(vv)

	return nil
}

// Value implements the driver Valuer interface.
// 将值转换为时间字符串
func (t Time) Value() (driver.Value, error) {
	return time.Time(t).Format(TimeFormat), nil
}

// UnmarshalJSON Json 转 Struct 数据结构
func (t *Time) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)
	return t.UnmarshalString(s)
}

// UnmarshalString string 转 Struct 数据结构
func (t *Time) UnmarshalString(data string) error {
	tm, err := time.ParseInLocation(TimeFormat, data, time.Local)
	if err != nil {
		return err
	}

	*t = Time(tm)
	return nil
}

// MarshalJSON Struct 转 Json 数据结构
// 格式成符合预期的时间格式字符串后返回，效果与Value一样
func (t Time) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("0000-00-00 00:00:00"), nil
	}
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')

	return b, nil
}

// GetDefaultPrepare 确定是否合法的default，不是交由后续处理
func (t Time) GetDefaultPrepare(rt reflect.Type, value string) (interface{}, error) {
	if value != "now" {
		return nil, &IgnoreError{}
	}
	return value, nil
}

// GetDefault 获取当前时间，只有这一种可能
func (t Time) GetDefault(context.Context, interface{}) interface{} {
	return Time(time.Now())
}
