package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/sourabh-kumar2/go-redis/store"
)

func TestEvalExpire(t *testing.T) {
	t.Parallel()

	s := store.New()

	cases := []struct {
		name    string
		setup   func(prefix string)
		args    func(prefix string) []string
		want    string
		wantErr bool
	}{
		{
			name:    "no args returns error",
			setup:   func(prefix string) {},
			args:    func(prefix string) []string { return []string{} },
			wantErr: true,
		},
		{
			name:    "only key no duration returns error",
			setup:   func(prefix string) {},
			args:    func(prefix string) []string { return []string{prefix + "k"} },
			wantErr: true,
		},
		{
			name:    "non-integer duration returns error",
			setup:   func(prefix string) {},
			args:    func(prefix string) []string { return []string{prefix + "k", "notanumber"} },
			wantErr: true,
		},
		{
			name:  "key does not exist returns 0",
			setup: func(prefix string) {},
			args:  func(prefix string) []string { return []string{prefix + "missing", "10"} },
			want:  ":0\r\n",
		},
		{
			name: "existing key gets expiry set and returns 1",
			setup: func(prefix string) {
				s.Put(prefix+"k", store.NewObj("val", -1))
			},
			args: func(prefix string) []string { return []string{prefix + "k", "10"} },
			want: ":1\r\n",
		},
		{
			name: "zero duration sets key to expire immediately",
			setup: func(prefix string) {
				s.Put(prefix+"k", store.NewObj("val", -1))
			},
			args: func(prefix string) []string { return []string{prefix + "k", "0"} },
			want: ":1\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prefix := "expire:" + t.Name() + ":"
			tc.setup(prefix)

			var buf bytes.Buffer
			err := evalExpire(tc.args(prefix), &buf, s)
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

func TestEvalExpireSetsExpiresAt(t *testing.T) {
	t.Parallel()

	s := store.New()
	key := "expire:expiresAt:" + t.Name()
	s.Put(key, store.NewObj("val", -1))

	before := time.Now().UnixMilli()
	var buf bytes.Buffer
	if err := evalExpire([]string{key, "5"}, &buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().UnixMilli()

	obj := s.Get(key)
	if obj == nil {
		t.Fatal("key unexpectedly missing after expire")
	}

	wantMin := before + 5000
	wantMax := after + 5000
	if obj.ExpiresAt < wantMin || obj.ExpiresAt > wantMax {
		t.Fatalf("ExpiresAt %d not in expected range [%d, %d]", obj.ExpiresAt, wantMin, wantMax)
	}
}
