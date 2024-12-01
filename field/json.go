package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value:%s", value)
	}

	if len(bytes) == 0 {
		bytes = []byte("null")
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	*j = JSON(result)

	return nil
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "null", nil
	}

	b, e := json.RawMessage(j).MarshalJSON()
	return string(b), e
}

// UnmarshalJSON Json 转 Struct 数据结构
func (j *JSON) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		var vv interface{}
		err := json.Unmarshal([]byte(v.(string)), &vv)
		if err != nil {
			return err
		}
		*j, _ = json.Marshal(vv)
	case reflect.Map:
		// 这个操作可以压缩传入数据的空白字符
		*j, _ = json.Marshal(v)
	case reflect.Slice:
		*j, _ = json.Marshal(v)
	default:
		return errors.New("type error, expect map []byte or json []byte received " + reflect.TypeOf(v).String())
	}

	return nil
}

// UnmarshalString string 转 Struct 数据结构
func (j *JSON) UnmarshalString(data string) error {
	return j.UnmarshalJSON([]byte(data))
}

// MarshalJSON Struct 转 Json 数据结构
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		j = []byte("null")
	}
	return j, nil
}
