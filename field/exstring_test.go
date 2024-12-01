package field

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExString(t *testing.T) {
	type S struct {
		Val ExString
	}

	tests := map[string]string{
		`{"val":null}`:  "null",
		`{"val":false}`: "false",
		`{"val":true}`:  "true",
		`{"val":1}`:     "1",
		`{"val":123.5}`: "123.5",
		`{"val":-1}`:    "-1",
		`{"val":0}`:     "0",
		`{"val":"1"}`:   "1",
		`{"val":""}`:    "",
		`{"val":"gdp"}`: "gdp",
		`{"val":"\"The string must be either double-quoted\"\n\u2639\u2639"}`: "\"The string must be either double-quoted\"\n\u2639\u2639",
	}
	for str, want := range tests {
		t.Run(str, func(t *testing.T) {
			s := &S{}
			if err := json.Unmarshal([]byte(str), s); err != nil {
				t.Fatal(err)
			}
			if s.Val.ToString() != want {
				t.Fatalf("s.Val=%q, want=%q", s.Val, want)
			}
		})
	}
}

func TestExString_MarshalJSON(t *testing.T) {
	type S1 struct {
		Val ExString
	}
	type S2 struct {
		Val string
	}
	tests := []struct {
		name   string
		fields *S1
		want   string
	}{
		{
			name: "case 1",
			fields: &S1{
				Val: "123",
			},
			want: `{"Val":"123"}`,
		},
		{
			name: "case 2",
			fields: &S1{
				Val: "123-你好",
			},
			want: `{"Val":"123-你好"}`,
		},
		{
			name: "case 3",
			fields: &S1{
				Val: "123-你好😆",
			},
			want: `{"Val":"123-你好😆"}`,
		},
		{
			name: "case 4",
			fields: &S1{
				Val: "我试试😘😘😘",
			},
			want: `{"Val":"我试试😘😘😘"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("ExString", func(t *testing.T) {
				bf, err := json.Marshal(tt.fields)
				require.NoError(t, err)
				require.Equal(t, tt.want, string(bf))

				var s1 *S1
				err2 := json.Unmarshal(bf, &s1)
				require.NoError(t, err2)
				require.NotNil(t, s1)
				require.Equal(t, tt.fields.Val, s1.Val)
			})

			t.Run("normal string Marshal s2", func(t *testing.T) {
				ss2 := &S2{
					Val: string(tt.fields.Val),
				}
				bf, err := json.Marshal(ss2)
				t.Logf("Marshal: %q", bf)
				require.NoError(t, err)

				var s1 *S1
				err2 := json.Unmarshal(bf, &s1)
				require.NoError(t, err2)
				require.Equal(t, tt.fields.Val, s1.Val)
			})
		})
	}
}

func TestExString1161(t *testing.T) {
	type FuzhenCard1 struct {
		Zhusu string `json:"zhusu,omitempty"`
	}
	type FuzhenCard2 struct {
		Zhusu ExString `json:"zhusu,omitempty"`
	}
	ss := `{"Zhusu":"\u6211\u8bd5\u8bd5\ud83d\ude18\ud83d\ude18\ud83d\ude18"}`
	var v1 *FuzhenCard1
	err := json.Unmarshal([]byte(ss), &v1)
	require.NoError(t, err)
	require.Equal(t, "我试试😘😘😘", v1.Zhusu)

	var v2 *FuzhenCard2
	err = json.Unmarshal([]byte(ss), &v2)
	require.NoError(t, err)
	require.Equal(t, "我试试😘😘😘", v2.Zhusu.ToString())

	bf1, err := json.Marshal(v2)
	require.NoError(t, err)

	bf2, err := json.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, string(bf1), string(bf2))
}

func TestExStringSlice(t *testing.T) {
	type user struct {
		Name string
		V1   *ExStringSlice
		V2   ExStringSlice
	}

	var v1 ExStringSlice
	user1 := &user{
		Name: "a",
		V1:   &v1,
	}

	checkUser := func(t *testing.T, got *user, wantV1 []string, wantV2 []string) {
		g1 := got.V1.StringSlice()
		if !reflect.DeepEqual(g1, wantV1) {
			t.Fatalf("wantV1=%v got=%v", wantV1, g1)
		}

		g2 := got.V2.StringSlice()
		if !reflect.DeepEqual(g2, wantV2) {
			t.Fatalf("wantV1=%v got=%v", wantV2, g2)
		}
	}

	var tests = []struct {
		name    string
		input   string
		want    *user
		wantErr bool
		check   func(t *testing.T, got *user)
	}{
		{
			name:    "case 1",
			input:   ``,
			wantErr: true,
		},
		{
			name:  "case 2",
			input: `{}`,
			want:  &user{},
		},
		{
			name:  "case 3",
			input: `{ "name":"a"}`,
			want: &user{
				Name: "a",
			},
		},
		{
			name:  "case 4-0",
			input: `{ "name":"a","v1":"abc","v2":["abc","def"]}`,
			want: &user{
				Name: "a",
				V1:   &ExStringSlice{"abc"},
				V2:   ExStringSlice{"abc", "def"},
			},
		},
		{
			name:  "case 4-1",
			input: `{ "name":"a","v1":"1,2","v2":["3","4"]}`,
			want: &user{
				Name: "a",
				V1:   &ExStringSlice{"1", "2"},
				V2:   ExStringSlice{"3", "4"},
			},
		},
		{
			name:  "case 4-2",
			input: `{ "name":"a","v1":"1,"}`,
			want: &user{
				Name: "a",
				V1:   &ExStringSlice{"1"},
			},
		},
		{
			name:  "case 5",
			input: `{ "name":"a","v1":"abc,hello, va,,","v2":"hello, vb,vc"}`,
			want: &user{
				Name: "a",
				V1:   &ExStringSlice{"abc", "hello", "va"},
				V2:   ExStringSlice{"hello", "vb", "vc"},
			},
			check: func(t *testing.T, got *user) {
				wantV1 := []string{"abc", "hello", "va"}
				wantV2 := []string{"hello", "vb", "vc"}
				checkUser(t, got, wantV1, wantV2)

				b, err := json.Marshal(got)
				if err != nil {
					t.Fatalf("err=%v", err)
				}
				want := []byte(`{"Name":"a","V1":["abc","hello","va"],"V2":["hello","vb","vc"]}`)
				if !bytes.Equal(b, want) {
					t.Fatalf("Marshal not Equal\n got=%q want=%q", b, want)
				}
			},
		},
		{
			name:  "case 6",
			input: `{ "name":"a","v1":null,"v2":null}`,
			want: &user{
				Name: "a",
			},
			check: func(t *testing.T, got *user) {
				checkUser(t, got, nil, nil)

				b, err := json.Marshal(got)
				if err != nil {
					t.Fatalf("err=%v", err)
				}
				want := []byte(`{"Name":"a","V1":null,"V2":null}`)
				if !bytes.Equal(b, want) {
					t.Fatalf("Marshal not Equal\n got=%q want=%q", b, want)
				}
			},
		},
		{
			name:    "case 7",
			input:   `{ "name":"a","v1":err,"v2":null}`, // 错误的格式
			wantErr: true,
		},
		{
			name:    "case 8",
			input:   `{ "name":"a","v1":{},"v2":null}`, // v1 不能是 object: {}
			wantErr: true,
			want:    user1,
		},
		{
			name:    "case 9",
			input:   `{ "name":"a","v1":{"a":"b"},"v2":null}`, // v1 不能是 object: {}
			wantErr: true,
			want:    user1,
		},
		{
			name:  "case 10",
			input: `{ "name":"a","v1":"","v2":[]}`,
			want:  user1,
		},
		{
			name:  "case 11",
			input: `{ "name":"a","v1":[],"v2":""}`,
			want:  user1,
			check: func(t *testing.T, got *user) {
				checkUser(t, got, nil, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *user
			err := json.Unmarshal([]byte(tt.input), &u)
			hasErr := err != nil
			if hasErr != tt.wantErr {
				t.Fatalf("wantErr=%v got=%v", tt.want, err)
			}
			if !reflect.DeepEqual(u, tt.want) {
				if u != nil {
					t.Logf("u.V1=%v isNil=%v, u.V2=%v isNil=%v", u.V1, u.V1 == nil, u.V2, u.V2 == nil)
				}
				t.Fatalf("not equal\n  got = %#v\n want = %#v", u, tt.want)
			}
			if tt.check != nil {
				tt.check(t, u)
			}
		})
	}
}

func TestExStringSlice_UnmarshalJSON(t *testing.T) {
	var s *ExStringSlice
	err := s.UnmarshalJSON(nil)
	if err == nil {
		t.Fatalf("expect has error")
	}
}
