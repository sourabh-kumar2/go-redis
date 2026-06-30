package core

import (
	"bytes"
	"testing"
)

func TestEvalDel(t *testing.T) {
	t.Parallel()

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
				Put(prefix+"k1", NewObj("v1", -1))
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
				Put(prefix+"a", NewObj("1", -1))
				Put(prefix+"b", NewObj("2", -1))
				Put(prefix+"c", NewObj("3", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "a", prefix + "b", prefix + "c"}
			},
			want: ":3\r\n",
		},
		{
			name: "delete mix of existing and missing keys counts only existing",
			setup: func(prefix string) {
				Put(prefix+"exists", NewObj("val", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "exists", prefix + "ghost"}
			},
			want: ":1\r\n",
		},
		{
			name: "delete same key twice counts only first deletion",
			setup: func(prefix string) {
				Put(prefix+"dup", NewObj("v", -1))
			},
			args: func(prefix string) []string {
				return []string{prefix + "dup", prefix + "dup"}
			},
			want: ":1\r\n",
		},
		{
			name: "deleted key is no longer retrievable",
			setup: func(prefix string) {
				Put(prefix+"gone", NewObj("v", -1))
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
			err := evalDel(tc.args(prefix), &buf)
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

	key := "del:remove:" + t.Name()
	Put(key, NewObj("val", -1))

	var buf bytes.Buffer
	if err := evalDel([]string{key}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if obj := Get(key); obj != nil {
		t.Fatalf("expected key %q to be deleted, but it still exists", key)
	}
}
