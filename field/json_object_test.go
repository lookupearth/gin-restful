package field

import (
	"encoding/json"
	"testing"
)

func TestJsonObject(t *testing.T) {
	jsonStr := `{"abc":1}`
	var j JSONObject
	if err := json.Unmarshal([]byte(jsonStr), &j); err != nil {
		t.Errorf("TestJsonObject Unmarshal fail, error=%v", err)
	}
	if j["abc"] != float64(1) {
		t.Errorf("TestJsonObject read key failed, expect=%v got=%v", float64(1), j["abc"])
	}
	v, err := json.Marshal(j)
	if err != nil {
		t.Errorf("TestJsonObject Marshal fail, error=%v", err)
	}
	if string(v) != jsonStr {
		t.Errorf("TestJsonObject Marshal fail, expect=%v got=%v", jsonStr, string(v))
	}

	jp := &JSONObject{}
	if err := jp.Scan([]byte(jsonStr)); err != nil {
		t.Errorf("TestJsonObject Scan fail, error=%v", err)
	}
	dbValue, err := jp.Value()
	if err != nil {
		t.Errorf("TestJsonObject Value fail, error=%v", err)
	}
	if dbValue != jsonStr {
		t.Errorf("TestJsonObject Marshal fail, expect=%v got=%v", jsonStr, dbValue)
	}

	j2 := &JSONObject{}
	if err := j2.UnmarshalString(jsonStr); err != nil {
		t.Errorf("TestJsonObject UnmarshalString fail, error=%v", err)
	}
	v2, err := json.Marshal(j)
	if err != nil {
		t.Errorf("TestJsonObject Marshal fail, error=%v", err)
	}
	if string(v2) != jsonStr {
		t.Errorf("TestJsonObject UnmarshalString fail, expect=%v got=%v", jsonStr, string(v2))
	}
}
