package core

import (
	"bytes"
	"testing"
	"time"
)

func TestEvalTTL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		setup   func(key string)
		want    string
		wantErr bool
	}{
		{
			name:  "missing key returns -2",
			setup: func(key string) {},
			want:  ":-2\r\n",
		},
		{
			name: "key with no expiry returns -1",
			setup: func(key string) {
				Put(key, NewObj("val", -1))
			},
			want: ":-1\r\n",
		},
		{
			name: "key with future expiry returns remaining seconds",
			setup: func(key string) {
				obj := &Obj{Value: "val", ExpiresAt: time.Now().UnixMilli() + 5000}
				Put(key, obj)
			},
			want: ":5\r\n",
		},
		{
			name: "already expired key returns -2",
			setup: func(key string) {
				obj := &Obj{Value: "val", ExpiresAt: time.Now().UnixMilli() - 1000}
				Put(key, obj)
			},
			want: ":-2\r\n",
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
			key := "ttl:" + t.Name()
			tc.setup(key)

			var args []string
			if tc.wantErr {
				args = []string{}
			} else {
				args = []string{key}
			}

			var buf bytes.Buffer
			err := evalTTL(args, &buf)
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
