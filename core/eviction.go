package core

// TODO: make the eviction strategy configurable
// TODO: Support multiple eviction strategy.
func evict() {
	evictFirst()
}

// Evicts the first key found while iterating the map
// TODO: use sampling to make it efficient
func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}
