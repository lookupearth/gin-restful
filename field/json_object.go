package field

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONObject map[string]interface{}

// Scan 读取db
func (j *JSONObject) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value:%s", value)
	}

	if len(bytes) == 0 {
		bytes = []byte("{}")
	}

	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	*j = result
	return nil
}

// Value 写入db
func (j JSONObject) Value() (driver.Value, error) {
	if j == nil {
		j = make(map[string]interface{})
	}
	b, e := json.Marshal(j)
	return string(b), e
}

// UnmarshalJSON 反序列化
func (j *JSONObject) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*j = v
	return nil
}

// UnmarshalString string 转 Struct 数据结构
func (j *JSONObject) UnmarshalString(data string) error {
	return j.UnmarshalJSON([]byte(data))
}

// MarshalJSON Struct 转 Json 数据结构
func (j JSONObject) MarshalJSON() ([]byte, error) {
	if j == nil {
		j = make(map[string]interface{})
	}
	return json.Marshal(map[string]interface{}(j))
}
