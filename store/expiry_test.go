package store

import (
	"testing"
	"time"
)

func TestExpireSample(t *testing.T) {
	cases := []struct {
		name        string
		setup       func(s *Store, prefix string)
		wantMinFrac float32
		wantMaxFrac float32
		wantRemoved func(prefix string) []string
	}{
		{
			name:        "no keys returns zero fraction",
			setup:       func(s *Store, prefix string) {},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
		{
			name: "all persistent keys returns zero fraction",
			setup: func(s *Store, prefix string) {
				for i := 0; i < 5; i++ {
					s.Put(prefix+string(rune('a'+i)), NewObj("v", -1))
				}
			},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
		{
			name: "all expired keys returns non-zero fraction",
			setup: func(s *Store, prefix string) {
				for i := 0; i < 5; i++ {
					obj := &Obj{Value: "v", ExpiresAt: time.Now().UnixMilli() - 1}
					s.Put(prefix+string(rune('a'+i)), obj)
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
			setup: func(s *Store, prefix string) {
				s.Put(prefix+"live", NewObj("v", 60000))
			},
			wantMinFrac: 0,
			wantMaxFrac: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			prefix := "expireSample:" + t.Name() + ":"
			tc.setup(s, prefix)

			frac := s.expireSample()

			if frac < tc.wantMinFrac || frac > tc.wantMaxFrac {
				t.Fatalf("fraction %v not in [%v, %v]", frac, tc.wantMinFrac, tc.wantMaxFrac)
			}

			if tc.wantRemoved != nil {
				for _, key := range tc.wantRemoved(prefix) {
					if s.Get(key) != nil {
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
		setup        func(s *Store, prefix string)
		wantDeleted  func(prefix string) []string
		wantSurvived func(prefix string) []string
	}{
		{
			name: "expired keys are removed",
			setup: func(s *Store, prefix string) {
				for i := 0; i < 5; i++ {
					obj := &Obj{Value: "v", ExpiresAt: time.Now().UnixMilli() - 1}
					s.Put(prefix+string(rune('a'+i)), obj)
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
			setup: func(s *Store, prefix string) {
				s.Put(prefix+"keep", NewObj("v", -1))
			},
			wantSurvived: func(prefix string) []string {
				return []string{prefix + "keep"}
			},
		},
		{
			name: "future-expiry keys survive",
			setup: func(s *Store, prefix string) {
				s.Put(prefix+"future", NewObj("v", 60000))
			},
			wantSurvived: func(prefix string) []string {
				return []string{prefix + "future"}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			prefix := "deleteExpired:" + t.Name() + ":"
			tc.setup(s, prefix)

			s.DeleteExpiredKeys()

			if tc.wantDeleted != nil {
				for _, key := range tc.wantDeleted(prefix) {
					if s.Get(key) != nil {
						t.Fatalf("expected key %q to be deleted but it still exists", key)
					}
				}
			}
			if tc.wantSurvived != nil {
				for _, key := range tc.wantSurvived(prefix) {
					if s.Get(key) == nil {
						t.Fatalf("expected key %q to survive but it was deleted", key)
					}
				}
			}
		})
	}
}
