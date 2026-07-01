package core

import (
	"log"
	"time"
)

func DeleteExpiredKeys() {
	for {
		frac := expireSample()
		if frac < 0.25 {
			break
		}
	}
	log.Println("deleted the expired but undeleted keys. total keys", StoreSize())
}

func expireSample() float32 {
	var limit int = 20
	var expiredCount int

	mu.Lock()
	defer mu.Unlock()

	for key, obj := range store {
		if obj.ExpiresAt == -1 {
			continue
		}

		limit--
		if obj.ExpiresAt <= time.Now().UnixMilli() {
			delete(store, key)
			expiredCount++
		}
		if limit == 0 {
			break
		}
	}
	return float32(expiredCount) / float32(20.0)
}
