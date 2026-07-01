package core

import (
	"bytes"
	"testing"

	"github.com/sourabh-kumar2/go-redis/store"
)

func TestEvalSET(t *testing.T) {
	t.Parallel()

	s := store.New()

	cases := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"missing value returns error", []string{"key"}, "", true},
		{"no args returns error", []string{}, "", true},
		{"key and value returns OK", []string{t.Name() + ":kv", "value"}, "+OK\r\n", false},
		{"key value with EX returns OK", []string{t.Name() + ":ex", "value", "EX", "10"}, "+OK\r\n", false},
		{"EX with no duration returns error", []string{t.Name() + ":noex", "value", "EX"}, "", true},
		{"EX with non-integer duration returns error", []string{t.Name() + ":badex", "value", "EX", "abc"}, "", true},
		{"unknown option returns error", []string{t.Name() + ":opt", "value", "UNKNOWN"}, "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := evalSET(tc.args, &buf, s)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
