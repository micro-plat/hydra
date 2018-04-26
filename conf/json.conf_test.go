package conf

import (
	"testing"

	"github.com/micro-plat/lib4go/ut"
)

func TestJsonConfNew(t *testing.T) {
	jc, err := NewJSONConf([]byte("{}"), 9999)
	ut.Expect(t, err, nil)
	ut.Expect(t, len(jc.data), 0)

	jc, err = NewJSONConf([]byte("{a}"), 9999)
	ut.Refute(t, err, nil)

	jc, err = NewJSONConf([]byte(`{"a":"b","b":100,"c":{"e":"f"}}`), 9999)
	ut.Expect(t, err, nil)
	ut.Expect(t, len(jc.data), 3)

	jc, err = NewJSONConf([]byte(`{"a":"b","b":100,"c":{"e":"f",a:100}}`), 9999)
	ut.Refute(t, err, nil)
}
func TestJsonConfVersion(t *testing.T) {
	jc, _ := NewJSONConf([]byte("{}"), 9999)
	ut.Expect(t, jc.GetVersion(), int32(9999))

	jc, _ = NewJSONConf([]byte("{}"), 0)
	ut.Expect(t, jc.GetVersion(), int32(0))

	jc, _ = NewJSONConf([]byte("{}"), -1)
	ut.Expect(t, jc.GetVersion(), int32(-1))

}
func TestJsonConfString(t *testing.T) {
	jc, _ := NewJSONConf([]byte("{}"), 9999)

	ut.Expect(t, jc.GetVersion(), int32(9999))
	ut.Expect(t, jc.GetString("a"), "")
	ut.Expect(t, jc.GetString("a", "b"), "b")

	jc, _ = NewJSONConf([]byte(`{"a":"b"}`), 9999)
	ut.Expect(t, jc.GetString("a"), "b")
	ut.Expect(t, jc.GetString("a", "c"), "b")

	jc, _ = NewJSONConf([]byte(`{"a":"100"}`), 9999)
	ut.Expect(t, jc.GetString("a"), "100")

	jc, _ = NewJSONConf([]byte(`{"a":true}`), 9999)
	ut.Expect(t, jc.GetString("a"), "true")

	jc, _ = NewJSONConf([]byte(`{"a":{"c":"d"}}`), 9999)
	ut.Refute(t, jc.GetString("a"), "b")

}
func TestJsonConfBool(t *testing.T) {
	jc, _ := NewJSONConf([]byte("{}"), 9999)

	b, err := jc.GetBool("a")
	ut.Refute(t, err, nil)
	ut.Expect(t, b, false)

	b, err = jc.GetBool("a", true)
	ut.Expect(t, err, nil)
	ut.Expect(t, b, true)
	tbs := []struct {
		input string
		value bool
	}{
		{input: `{"a":"T"}`, value: true},
		{input: `{"a":"t"}`, value: true},
		{input: `{"a":"true"}`, value: true},
		{input: `{"a":true}`, value: true},
		{input: `{"a":"TRUE"}`, value: true},
		{input: `{"a":"1"}`, value: true},
		{input: `{"a":"Y"}`, value: true},
		{input: `{"a":"yes"}`, value: true},
		{input: `{"a":"on"}`, value: true},
		{input: `{"a":"no"}`, value: false},
		{input: `{"a":"0"}`, value: false},
		{input: `{"a":"abc"}`, value: false},
	}

	for _, tb := range tbs {
		jc, _ = NewJSONConf([]byte(tb.input), 9999)

		b, _ = jc.GetBool("a")
		ut.Expect(t, b, tb.value)
	}

}

func TestJsonConfInt(t *testing.T) {
	tbs :=
		[]struct {
			input    string
			field    string
			version  int32
			def      int
			expected int
		}{
			{input: "{}", field: "a", expected: 0},
			{input: "{}", field: "a", def: 2, expected: 2},
			{input: `{"a":"c"}`, field: "a", def: 3, expected: 3},
			{input: `{"a":"c,b"}`, field: "a", def: 4, expected: 4},
			{input: `{"a":"c;b"}`, field: "a", def: 5, expected: 5},
			{input: `{"a":"c/b"}`, field: "a", def: 6, expected: 6},
			{input: `{"a":100}`, field: "a", def: 6, expected: 100},
			{input: `{"a":-100}`, field: "a", def: 6, expected: -100},
			{input: `{"a":2147483646}`, field: "a", def: 6, expected: 6},
		}

	for _, tb := range tbs {
		jc, _ := NewJSONConf([]byte(tb.input), tb.version)
		ut.Expect(t, jc.GetInt(tb.field, tb.def), tb.expected)
	}
}

func TestJsonconfStrings(t *testing.T) {
	tbs :=
		[]struct {
			input    string
			field    string
			def      []string
			len      int
			expected []string
		}{
			{input: "{}", field: "a", len: 0},
			{input: "{}", field: "a", def: []string{"a", "b"}, len: 2, expected: []string{"a", "b"}},
			{input: `{"a":"c"}`, field: "a", def: []string{"a", "b"}, len: 1, expected: []string{"c"}},
			{input: `{"a":"c,b"}`, field: "a", def: []string{"a", "b"}, len: 1, expected: []string{"c,b"}},
			{input: `{"a":"c;b"}`, field: "a", def: []string{"a"}, len: 2, expected: []string{"c", "b"}},
			{input: `{"a":"c/b"}`, field: "a", def: []string{"a"}, len: 1, expected: []string{"c/b"}},
		}

	for _, tb := range tbs {
		jc, _ := NewJSONConf([]byte(tb.input), 9999)
		value := jc.GetStrings(tb.field, tb.def...)
		ut.ExpectSkip(t, len(value), tb.len)
		for i, v := range tb.expected {
			ut.Expect(t, value[i], v)
		}
	}
}
func TestJsonConfSecion(t *testing.T) {
	jc, _ := NewJSONConf([]byte(`{"a":"b","b":{"c":100}}`), 9999)

	b, err := jc.GetSection("a")
	ut.Refute(t, err, nil)

	b, err = jc.GetSection("b")
	ut.Expect(t, err, nil)
	ut.Expect(t, b.version, int32(9999))
	ut.Expect(t, b.GetString("c"), "100")
	ut.Expect(t, b.GetInt("c"), 100)
}
