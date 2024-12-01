package field

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestExInt64(t *testing.T) {
	type S struct {
		Val ExInt64
	}
	cases := map[string]int64{
		`{"val":12345}`:   12345,
		`{"val":"12345"}`: 12345,

		`{"val":-12345}`:   -12345,
		`{"val":"-12345"}`: -12345,

		`{"val":"12345.1"}`:    12345,
		`{"val":12345.1}`:      12345,
		`{"val":"12345.0"}`:    12345,
		`{"val":"12345.0001"}`: 12345,

		`{"val":"123456789012345678"}`:  123456789012345678,
		`{"val":"-123456789012345678"}`: -123456789012345678,

		`{"val":"-12345.1"}`:    -12345,
		`{"val":-12345.1}`:      -12345,
		`{"val":"-12345.0"}`:    -12345,
		`{"val":"-12345.0001"}`: -12345,

		`{"val":false}`: 0,
		`{"val":true}`:  1,
		`{"val":null}`:  0,
		`{"val":""}`:    0,
	}

	for caseData, expectVal := range cases {
		s := &S{}
		if err := json.Unmarshal([]byte(caseData), s); err != nil {
			t.Errorf("json.Unmarshal %s failed: %s", caseData, err.Error())
		}
		if s.Val.ToInt64() != expectVal {
			t.Errorf("case = %s : expect = %d,real= %d", caseData, expectVal, s.Val)
		}
	}

	notAllowCases := []string{
		`{"Val":"other"}`,
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

func BenchmarkExInt64_UnmarshalJSON1(b *testing.B) {
	type S struct {
		Val ExInt64
	}
	data := []byte(`{"Val":1}`)
	b.ResetTimer()
	var s *S
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, s)
	}
}

func BenchmarkExInt64_UnmarshalJSON2(b *testing.B) {
	type S struct {
		Val ExInt64
	}
	data := []byte(`{"Val":false}`)
	b.ResetTimer()
	var s *S
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, s)
	}
}

func TestExInt64_MarshalJSON(t *testing.T) {
	type S struct {
		Val ExInt64
	}
	s := &S{
		Val: 100,
	}
	bf, err := json.Marshal(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	want := []byte(`{"Val":100}`)
	if !bytes.Equal(bf, want) {
		t.Fatalf("got=%q want=%q", bf, want)
	}
}
