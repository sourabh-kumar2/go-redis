package store

import (
	"log"
	"time"
)

func (s *Store) DeleteExpiredKeys() {
	for {
		frac := s.expireSample()
		if frac < 0.25 {
			break
		}
	}
	log.Println("deleted the expired but undeleted keys. total keys", s.Size())
}

func (s *Store) expireSample() float32 {
	var limit int = 20
	var expiredCount int

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, obj := range s.items {
		if obj.ExpiresAt == -1 {
			continue
		}

		limit--
		if obj.ExpiresAt <= time.Now().UnixMilli() {
			delete(s.items, key)
			expiredCount++
		}
		if limit == 0 {
			break
		}
	}
	return float32(expiredCount) / float32(20.0)
}
