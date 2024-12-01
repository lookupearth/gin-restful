package field

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	timeStr := `"2022-06-28 13:49:05"`
	timeClear := `2022-06-28 13:49:05`
	var field Time
	if err := json.Unmarshal([]byte(timeStr), &field); err != nil {
		t.Errorf("TestTime Unmarshal fail, error=%v", err)
	}
	v, err := json.Marshal(field)
	if err != nil {
		t.Errorf("TestTime Marshal fail, error=%v", err)
	}
	if string(v) != timeStr {
		t.Errorf("TestJson Marshal fail, expect=%v got=%v", timeStr, string(v))
	}

	tm := time.Date(2022, 06, 28, 13, 49, 05, 00, time.Local)
	fieldPtr := &Time{}
	if err := fieldPtr.Scan(tm); err != nil {
		t.Errorf("TestTime Scan fail, error=%v", err)
	}
	dbValue, err := fieldPtr.Value()
	if err != nil {
		t.Errorf("TestTime Value fail, error=%v", err)
	}
	if dbValue != timeClear {
		t.Errorf("TestTime Marshal fail, expect=%v got=%v", timeClear, dbValue)
	}
	v2, err := json.Marshal(fieldPtr)
	if err != nil {
		t.Errorf("TestTime Marshal fail, error=%v", err)
	}
	if string(v2) != timeStr {
		t.Errorf("TestJson Marshal fail, expect=%v got=%v", timeStr, string(v2))
	}

	field2 := &Time{}
	if err := field2.UnmarshalString(timeClear); err != nil {
		t.Errorf("TestTime UnmarshalJSON fail, error=%v", err)
	}
	if time.Time(*field2) != tm {
		t.Errorf("TestTime UnmarshalJSON fail, expect=%v got=%v", tm, field2)
	}

}
