package core

import (
	"sync"
	"time"
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
	store[k] = obj
}

func Get(k string) *Obj {
	mu.RLock()
	defer mu.RUnlock()
	return store[k]
}
