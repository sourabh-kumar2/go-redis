package core

import (
	"testing"
	"time"
)

func TestExpireSample(t *testing.T) {
	cases := []struct {
		name        string
		setup       func(prefix string)
		wantMinFrac float32
		wantMaxFrac float32
		wantRemoved func(prefix string) []string
	}{
		{
			name:        "no keys returns zero fraction",
			setup:       func(prefix string) {},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
		{
			name: "all persistent keys returns zero fraction",
			setup: func(prefix string) {
				for i := 0; i < 5; i++ {
					Put(prefix+string(rune('a'+i)), NewObj("v", -1))
				}
			},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
		{
			name: "all expired keys returns non-zero fraction",
			setup: func(prefix string) {
				for i := 0; i < 5; i++ {
					obj := &Obj{Value: "v", ExpiresAt: time.Now().UnixMilli() - 1}
					Put(prefix+string(rune('a'+i)), obj)
				}
			},
			wantMinFrac: 0.01,
			wantMaxFrac: 1.0,
			wantRemoved: func(prefix string) []string {
				keys := make([]string, 5)
				for i := range keys {
					keys[i] = prefix + string(rune('a'+i))
				}
				return keys
			},
		},
		{
			name: "future-expiry keys are not deleted",
			setup: func(prefix string) {
				Put(prefix+"live", NewObj("v", 60000))
			},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefix := "expireSample:" + t.Name() + ":"
			tc.setup(prefix)

			frac := expireSample()

			if frac < tc.wantMinFrac || frac > tc.wantMaxFrac {
				t.Fatalf("fraction %v not in [%v, %v]", frac, tc.wantMinFrac, tc.wantMaxFrac)
			}

			if tc.wantRemoved != nil {
				for _, key := range tc.wantRemoved(prefix) {
					if Get(key) != nil {
						t.Fatalf("expected key %q to be deleted but it still exists", key)
					}
				}
			}
		})
	}
}

func TestDeleteExpiredKeys(t *testing.T) {
	cases := []struct {
		name         string
		setup        func(prefix string)
		wantDeleted  func(prefix string) []string
		wantSurvived func(prefix string) []string
	}{
		{
			name: "expired keys are removed",
			setup: func(prefix string) {
				for i := 0; i < 5; i++ {
					obj := &Obj{Value: "v", ExpiresAt: time.Now().UnixMilli() - 1}
					Put(prefix+string(rune('a'+i)), obj)
				}
			},
			wantDeleted: func(prefix string) []string {
				keys := make([]string, 5)
				for i := range keys {
					keys[i] = prefix + string(rune('a'+i))
				}
				return keys
			},
		},
		{
			name: "persistent keys survive",
			setup: func(prefix string) {
				Put(prefix+"keep", NewObj("v", -1))
			},
			wantSurvived: func(prefix string) []string {
				return []string{prefix + "keep"}
			},
		},
		{
			name: "future-expiry keys survive",
			setup: func(prefix string) {
				Put(prefix+"future", NewObj("v", 60000))
			},
			wantSurvived: func(prefix string) []string {
				return []string{prefix + "future"}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefix := "deleteExpired:" + t.Name() + ":"
			tc.setup(prefix)

			DeleteExpiredKeys()

			if tc.wantDeleted != nil {
				for _, key := range tc.wantDeleted(prefix) {
					if Get(key) != nil {
						t.Fatalf("expected key %q to be deleted but it still exists", key)
					}
				}
			}
			if tc.wantSurvived != nil {
				for _, key := range tc.wantSurvived(prefix) {
					if Get(key) == nil {
						t.Fatalf("expected key %q to survive but it was deleted", key)
					}
				}
			}
		})
	}
}
