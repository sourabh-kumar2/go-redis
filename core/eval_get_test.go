package core

import (
	"bytes"
	"testing"
	"time"
)

func TestEvalGET(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		setup   func(key string)
		want    string
		wantErr bool
	}{
		{
			name:  "missing key returns nil",
			setup: func(key string) {},
			want:  RESP_NIL,
		},
		{
			name: "existing key returns bulk string",
			setup: func(key string) {
				Put(key, NewObj("world", -1))
			},
			want: "$5\r\nworld\r\n",
		},
		{
			name: "expired key returns nil",
			setup: func(key string) {
				obj := NewObj("ghost", 1)
				obj.ExpiresAt = time.Now().UnixMilli() - 1
				Put(key, obj)
			},
			want: RESP_NIL,
		},
		{
			name: "key with future expiry returns value",
			setup: func(key string) {
				Put(key, NewObj("alive", 60000))
			},
			want: "$5\r\nalive\r\n",
		},
		{
			name:    "wrong number of args returns error",
			setup:   func(key string) {},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			key := "get:" + t.Name()
			tc.setup(key)

			var args []string
			if tc.wantErr {
				args = []string{} // zero args triggers the error
			} else {
				args = []string{key}
			}

			var buf bytes.Buffer
			err := evalGET(args, &buf)
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
