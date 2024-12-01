package field

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

var _ json.Marshaler = (*ExString)(nil)
var _ json.Unmarshaler = (*ExString)(nil)

// ExString 字符串类型
//
//	支持值形式如："abcde", 12345, true, false, null
//	若传入为 bool 类型，会获的其字符串形式的值："true"、"false"
//	若传入的是 null，会获得其字符串形式的值："null"
//	若传入的是数字，如 123,会获的其字符串形式的值: "123"
type ExString string

// UnmarshalJSON for ExString
func (a *ExString) UnmarshalJSON(b []byte) (e error) {
	// 0x22 是 "
	if len(b) > 1 && b[0] == 0x22 && b[len(b)-1] == 0x22 {
		var s1 string
		if err := json.Unmarshal(b, &s1); err != nil {
			return err
		}
		*a = ExString(s1)
		return nil
	}
	s := string(b)
	*a = ExString(s)
	return nil
}

// MarshalJSON  编码为 json
func (a ExString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(a))
}

// ToString underlying type
func (a ExString) ToString() string {
	return string(a)
}

var _ json.Unmarshaler = (*ExStringSlice)(nil)

// ExStringSlice 扩展 json 数组类型，
// 允许格式："abc"、["avc","def"]、"abc,def"、null、"1,2"、["1","2"]
//
// 若字段定义为该类型，在 json.Marshal 的时候，当值为空的时候会编码为 null。
type ExStringSlice []string

// UnmarshalJSON implement json interface
func (s *ExStringSlice) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// 正常情况下，使用 json.Unmarshal 是不会运行到这个逻辑的
		return errEmptyInput
	}

	if bytes.Equal(b, null) {
		return nil
	}

	switch b[0] {
	case '[':
		var arr []string
		err := json.Unmarshal(b, &arr)
		if err == nil && len(arr) > 0 {
			*s = arr
		}
		return err
	case '"':
		var vs string
		err := json.Unmarshal(b, &vs)
		if err == nil && len(vs) > 0 {
			arr := strings.Split(vs, ",")
			for i := 0; i < len(arr); i++ {
				item := strings.TrimSpace(arr[i])
				if len(item) > 0 {
					*s = append(*s, item)
				}
			}
		}
		return err
	default:
		return fmt.Errorf("cannnot parser as ExStringSlice, first char is %c", b[0])
	}
}

// StringSlice 返回 []string 的值
func (s *ExStringSlice) StringSlice() []string {
	if s == nil {
		return nil
	}
	return *s
}
