package store

// TODO: make the eviction strategy configurable
// TODO: Support multiple eviction strategies.
func (s *Store) evict() {
	s.evictFirst()
}

// evictFirst evicts the first key found while iterating the map.
// Must be called with s.mu.Lock already held.
// TODO: use sampling to make it efficient
func (s *Store) evictFirst() {
	for k := range s.items {
		delete(s.items, k)
		return
	}
}
