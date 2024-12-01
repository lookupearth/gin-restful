package field

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	jsonStr := `[{"abc":1}]`
	var j JSON
	if err := json.Unmarshal([]byte(jsonStr), &j); err != nil {
		t.Errorf("TestJson Unmarshal fail, error=%v", err)
	}
	v, err := json.Marshal(j)
	if err != nil {
		t.Errorf("TestJson Marshal fail, error=%v", err)
	}
	if string(v) != jsonStr {
		t.Errorf("TestJson Marshal fail, expect=%v got=%v", jsonStr, string(v))
	}

	jp := &JSON{}
	if err := jp.Scan([]byte(jsonStr)); err != nil {
		t.Errorf("TestJson Scan fail, error=%v", err)
	}
	dbValue, err := jp.Value()
	if err != nil {
		t.Errorf("TestJson Value fail, error=%v", err)
	}
	if dbValue != jsonStr {
		t.Errorf("TestJson Marshal fail, expect=%v got=%v", jsonStr, dbValue)
	}

	j2 := &JSON{}
	if err := j2.UnmarshalString(jsonStr); err != nil {
		t.Errorf("TestJson UnmarshalString fail, error=%v", err)
	}
	v2, err := json.Marshal(j)
	if err != nil {
		t.Errorf("TestJson Marshal fail, error=%v", err)
	}
	if string(v2) != jsonStr {
		t.Errorf("TestJson UnmarshalString fail, expect=%v got=%v", jsonStr, string(v2))
	}
}
