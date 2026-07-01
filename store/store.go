package store

import (
	"sync"
	"time"

	"github.com/sourabh-kumar2/go-redis/config"
)

type Store struct {
	mu    sync.RWMutex
	items map[string]*Obj
}

type Obj struct {
	Value     any
	ExpiresAt int64
}

func New() *Store {
	return &Store{items: make(map[string]*Obj)}
}

func NewObj(value any, durationMs int64) *Obj {
	expiresAt := int64(-1)
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	return &Obj{Value: value, ExpiresAt: expiresAt}
}

func (s *Store) Put(k string, obj *Obj) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) >= config.KeysLimit {
		s.evict()
	}
	s.items[k] = obj
}

func (s *Store) Get(k string) *Obj {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items[k]
}

func (s *Store) Del(k string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[k]; ok {
		delete(s.items, k)
		return true
	}
	return false
}

func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
