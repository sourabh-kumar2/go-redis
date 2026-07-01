package core

import (
	"sync"
	"time"

	"github.com/sourabh-kumar2/go-redis/config"
)

var (
	mu    sync.RWMutex
	store map[string]*Obj
)

func init() {
	store = make(map[string]*Obj)
}

type Obj struct {
	Value     any
	ExpiresAt int64
}

func NewObj(value any, durationMs int64) *Obj {
	expiresAt := int64(-1)
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}

	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, obj *Obj) {
	mu.Lock()
	defer mu.Unlock()
	// evict is called with the lock already held — evictFirst must not re-acquire it.
	if len(store) >= config.KeysLimit {
		evict()
	}
	store[k] = obj
}

func Get(k string) *Obj {
	mu.RLock()
	defer mu.RUnlock()
	return store[k]
}

func Del(k string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := store[k]; ok {
		delete(store, k)
		return true
	}
	return false
}

func StoreSize() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(store)
}
