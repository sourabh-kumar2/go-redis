package core

import (
	"bytes"
	"testing"
)

func TestEvalPING(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"no args returns PONG", []string{}, "+PONG\r\n", false},
		{"one arg echoes as bulk string", []string{"hello"}, "$5\r\nhello\r\n", false},
		{"two args returns error", []string{"a", "b"}, "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := evalPING(tc.args, &buf)
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
