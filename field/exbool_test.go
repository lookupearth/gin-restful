package field

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestExBool(t *testing.T) {
	type S struct {
		Val ExBool
	}

	cases := map[string]bool{
		`{"Val":false}`:   false,
		`{"Val":"false"}`: false,
		`{"Val":"f"}`:     false,
		`{"Val":"FALSE"}`: false,
		`{"Val":"False"}`: false,
		`{"Val":"0"}`:     false,
		`{"Val":0}`:       false,
		`{"Val":true}`:    true,
		`{"Val":"true"}`:  true,
		`{"Val":"True"}`:  true,
		`{"Val":"TRUE"}`:  true,
		`{"Val":1}`:       true,
		`{"Val":"1"}`:     true,
		`{"Val":"t"}`:     true,
		`{"Val":"T"}`:     true,
		`{}`:              false,
		`{"Val":null}`:    false,
		`{"Val":"null"}`:  false,
	}
	trueBf := []byte(`{"Val":true}`)
	falseBf := []byte(`{"Val":false}`)
	for str, want := range cases {
		t.Run(str, func(t *testing.T) {
			s := &S{}
			if err := json.Unmarshal([]byte(str), s); err != nil {
				t.Errorf("json.Unmarshalfailed: %v", err)
			}
			if got := s.Val.ToBool(); got != want {
				t.Fatalf("ToBool()=%v want=%v", got, want)
			}
			t.Run("Marshal", func(t *testing.T) {
				bf, err := json.Marshal(s)
				if err != nil {
					t.Fatal(err.Error())
				}
				var wantBt []byte
				if want {
					wantBt = trueBf
				} else {
					wantBt = falseBf
				}
				if !bytes.Equal(bf, wantBt) {
					t.Fatalf("got=%q want=%q", bf, wantBt)
				}
			})
		})
	}
	notAllowCases := []string{
		`{"Val":"2"}`,
		`{"Val":2}`,
	}

	for _, str := range notAllowCases {
		t.Run(str, func(t *testing.T) {
			s := &S{}
			if err := json.Unmarshal([]byte(str), s); err == nil {
				t.Fatal("expect has error")
			}
		})
	}
}
