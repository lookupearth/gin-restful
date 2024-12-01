package field

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

var _ json.Marshaler = (*ExInt64)(nil)
var _ json.Unmarshaler = (*ExInt64)(nil)

// ExInt64 扩展的 int64 类型
//
// 支持值形式如: 12345, "12345", false(解析为 0), null(解析为 0), "12345.1", "12345.0", 12345.1
//
//	"" 空字符串解析为 0
//	true 解析为 1
type ExInt64 int64

// UnmarshalJSON for ExInt64
func (a *ExInt64) UnmarshalJSON(b []byte) error {
	str := string(bytes.Trim(b, `"`))
	strParts := strings.Split(str, ".")
	switch len(strParts) {
	case 2:
		// 简单判断是否为 float 类型（逻辑不严谨）
		fallthrough
	case 1:
		// int 类型
		n, err := strconv.ParseInt(strParts[0], 10, 64)
		if err == nil {
			*a = ExInt64(n)
			return nil
		}
	}
	switch str {
	case "false", "null", "":
		*a = 0
		return nil
	case "true":
		*a = 1
		return nil
	}
	return &json.InvalidUnmarshalError{Type: reflect.TypeOf(a)}
}

// MarshalJSON to 12345
func (a ExInt64) MarshalJSON() ([]byte, error) {
	s := strconv.FormatInt(int64(a), 10)
	return []byte(s), nil
}

// ToInt64 获取值
func (a ExInt64) ToInt64() int64 {
	return int64(a)
}
