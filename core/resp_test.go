package core

import (
	"fmt"
	"testing"
)

func TestSimpleStringDecode(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"+OK\r\n": "OK",
	}
	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fatalf("got: %v, want: %v", value, v)
		}
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"-Error Message\r\n": "Error Message",
	}
	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fatalf("got: %v, want: %v", value, v)
		}
	}
}

func TestInt64Decode(t *testing.T) {
	t.Parallel()

	cases := map[string]int64{
		":0\r\n":    0,
		":1000\r\n": 1000,
		":-100\r\n": -100,
	}
	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fatalf("got: %v, want: %v", value, v)
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}
	for k, v := range cases {
		value, _ := Decode([]byte(k))
		if v != value {
			t.Fatalf("got: %v, want: %v", value, v)
		}
	}
}

func TestArrayDecode(t *testing.T) {
	t.Parallel()

	cases := map[string][]any{
		"*0\r\n":                                        {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":          {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                      {1, 2, 3},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n": {1, 2, 3, 4, "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n$5\r\nhello\r\n$5\r\nWorld\r\n": {[]int64{1, 2, 3}, []string{"hello", "World"}},
	}
	for k, v := range cases {
		value, _ := Decode([]byte(k))

		array := value.([]any)
		if len(array) != len(v) {
			t.Fatalf("got length: %v, want length: %v", len(array), len(v))
		}
		for i := range array {
			if fmt.Sprintf("%v", array[i]) != fmt.Sprintf("%v", v[i]) {
				t.Fatalf("got: %v, want: %v", array[i], v[i])
			}
		}
	}
}
