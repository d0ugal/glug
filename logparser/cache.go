package logparser

import (
	"sync"
)

// timestampFieldCache caches the results of isTimestampField checks
type timestampFieldCache struct {
	mu    sync.RWMutex
	cache map[string]bool
}

var fieldCache = &timestampFieldCache{
	cache: make(map[string]bool),
}

// isTimestampFieldCached checks if a field is a timestamp field with caching
func isTimestampFieldCached(fieldName string) bool {
	fieldCache.mu.RLock()
	if cached, exists := fieldCache.cache[fieldName]; exists {
		fieldCache.mu.RUnlock()
		return cached
	}
	fieldCache.mu.RUnlock()

	// Not in cache, compute and store
	result := isTimestampField(fieldName)
	
	fieldCache.mu.Lock()
	fieldCache.cache[fieldName] = result
	fieldCache.mu.Unlock()
	
	return result
}

// clearTimestampFieldCache clears the cache (useful for testing)
func clearTimestampFieldCache() {
	fieldCache.mu.Lock()
	fieldCache.cache = make(map[string]bool)
	fieldCache.mu.Unlock()
}
