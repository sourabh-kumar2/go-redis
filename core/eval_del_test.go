package core

import (
	"bytes"
	"testing"

	"github.com/sourabh-kumar2/go-redis/store"
)

func TestEvalDel(t *testing.T) {
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
			name: "delete existing key returns 1",
			setup: func(prefix string) {
				s.Put(prefix+"k1", store.NewObj("v1", -1))
			},
			args: func(prefix string) []string { return []string{prefix + "k1"} },
			want: ":1\r\n",
		},
		{
			name:  "delete missing key returns 0",
			setup: func(prefix string) {},
			args:  func(prefix string) []string { return []string{prefix + "missing"} },
			want:  ":0\r\n",
		},
		{
			name: "delete multiple existing keys returns count",
			setup: func(prefix string) {
				s.Put(prefix+"a", store.NewObj("1", -1))
				s.Put(prefix+"b", store.NewObj("2", -1))
				s.Put(prefix+"c", store.NewObj("3", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "a", prefix + "b", prefix + "c"}
			},
			want: ":3\r\n",
		},
		{
			name: "delete mix of existing and missing keys counts only existing",
			setup: func(prefix string) {
				s.Put(prefix+"exists", store.NewObj("val", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "exists", prefix + "ghost"}
			},
			want: ":1\r\n",
		},
		{
			name: "delete same key twice counts only first deletion",
			setup: func(prefix string) {
				s.Put(prefix+"dup", store.NewObj("v", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "dup", prefix + "dup"}
			},
			want: ":1\r\n",
		},
		{
			name: "deleted key is no longer retrievable",
			setup: func(prefix string) {
				s.Put(prefix+"gone", store.NewObj("v", -1))
			},
			args: func(prefix string) []string { return []string{prefix + "gone"} },
			want: ":1\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prefix := "del:" + t.Name() + ":"
			tc.setup(prefix)

			var buf bytes.Buffer
			err := evalDel(tc.args(prefix), &buf, s)
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

func TestEvalDelRemovesKey(t *testing.T) {
	t.Parallel()

	s := store.New()
	key := "del:remove:" + t.Name()
	s.Put(key, store.NewObj("val", -1))

	var buf bytes.Buffer
	if err := evalDel([]string{key}, &buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if obj := s.Get(key); obj != nil {
		t.Fatalf("expected key %q to be deleted, but it still exists", key)
	}
}
