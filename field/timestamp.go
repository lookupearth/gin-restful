package field

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Timestamp time.Time

// Scan implements the Scanner interface.
// 将原始数据进行Struct结构改造
func (t *Timestamp) Scan(value interface{}) error {
	// set default value
	*t = Timestamp(time.Time{})

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

	*t = Timestamp(vv)

	return nil
}

// Value implements the driver Valuer interface.
// 将值转换为时间字符串
func (t Timestamp) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return "0000-00-00 00:00:00", nil
	}
	return time.Time(t).Format(TimeFormat), nil
}

// UnmarshalJSON Json 转 Struct 数据结构
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	return t.UnmarshalString(string(data))
}

// UnmarshalString string 转 Struct 数据结构
func (t *Timestamp) UnmarshalString(data string) error {
	timeUnix, err := strconv.ParseInt(data, 0, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %s for timestamp, got error: %v", string(data), err)
	}

	*t = Timestamp(time.Unix(timeUnix, 0))
	return nil
}

// MarshalJSON Struct 转 Json 数据结构
// 格式成符合预期的时间格式字符串后返回，效果与Value一样
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("0"), nil
	}
	ts := time.Time(t).Unix()
	b := strconv.FormatInt(ts, 10)
	return []byte(b), nil
}

// GetDefaultPrepare 确定是否合法的default，不是交由后续处理
func (t Timestamp) GetDefaultPrepare(rt reflect.Type, value string) (interface{}, error) {
	if value != "now" {
		return nil, &IgnoreError{}
	}
	return value, nil
}

// GetDefault 获取当前时间，只有这一种可能
func (t Timestamp) GetDefault(context.Context, interface{}) interface{} {
	return Timestamp(time.Now())
}
