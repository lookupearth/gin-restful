package field

import (
	"bytes"
	"encoding/json"
	"strconv"
)

var _ json.Marshaler = (*ExBool)(nil)
var _ json.Unmarshaler = (*ExBool)(nil)

// ExBool bool 类型
//
// 支持值形式如: "true", 1, true, false, null
//
//	若是数字，只支持 0( 包含 "0") 和 1 ( 包含 "1")，若传入数字 2 或者其他值，将报错
type ExBool bool

// UnmarshalJSON for ExString,
func (a *ExBool) UnmarshalJSON(b []byte) error {
	str := string(bytes.Trim(b, `"`))
	if str == "null" {
		*a = false
		return nil
	}
	r, e := strconv.ParseBool(str)
	if e != nil {
		return e
	}
	*a = ExBool(r)
	return nil
}

// MarshalJSON to "true", "false"
func (a ExBool) MarshalJSON() ([]byte, error) {
	if a {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// ToBool 读取值
func (a ExBool) ToBool() bool {
	return bool(a)
}
