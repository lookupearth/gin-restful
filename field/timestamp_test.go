package field

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	timeClear := `2022-06-28 13:49:05`
	ts := "1656395345"
	var field Timestamp
	if err := json.Unmarshal([]byte(ts), &field); err != nil {
		t.Errorf("TestTimestamp Unmarshal fail, error=%v", err)
	}
	v, err := json.Marshal(field)
	if err != nil {
		t.Errorf("TestTimestamp Marshal fail, error=%v", err)
	}
	if string(v) != ts {
		t.Errorf("TestTimestamp Marshal fail, expect=%v got=%v", ts, string(v))
	}

	tm := time.Date(2022, 06, 28, 13, 49, 05, 00, time.Local)
	fieldPtr := &Timestamp{}
	if err := fieldPtr.Scan(tm); err != nil {
		t.Errorf("TestTimestamp Scan fail, error=%v", err)
	}
	dbValue, err := fieldPtr.Value()
	if err != nil {
		t.Errorf("TestTimestamp Value fail, error=%v", err)
	}
	if dbValue != timeClear {
		t.Errorf("TestTimestamp Marshal fail, expect=%v got=%v", timeClear, dbValue)
	}
	v2, err := json.Marshal(fieldPtr)
	if err != nil {
		t.Errorf("TestTimestamp Marshal fail, error=%v", err)
	}
	if string(v2) != ts {
		t.Errorf("TestTimestamp Marshal fail, expect=%v got=%v", ts, string(v2))
	}

	field2 := &Timestamp{}
	if err := field2.UnmarshalString(ts); err != nil {
		t.Errorf("TestTimestamp UnmarshalString fail, error=%v", err)
	}
	if time.Time(*field2) != tm {
		t.Errorf("TestTimestamp UnmarshalString fail, expect=%v got=%v", tm, time.Time(*field2))
	}
}
