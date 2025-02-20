package cache

import (
	"sync"
	"time"
)

var instance *CacheService

type (
	// Implements Global value cache for web-server
	CacheService struct {
		store ProvideMap
		mu    sync.RWMutex
	}

	ProvideMap map[any]CacheDataUnit

	CacheDataUnit struct {
		CreatedAt time.Time
		CachedObj any

		_ struct{}
	}
)

func Init() {
	instance = &CacheService{
		store: ProvideMap{},
	}
}

func GetEntry() *CacheService {
	return instance
}

func (cache *CacheService) GetValue(key any) any {

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	v, ok := cache.store[key]
	if ok {
		return v.CachedObj
	}
	return nil
}

func (cache *CacheService) PushValue(key, value any) {

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.store[key] = CacheDataUnit{
		CreatedAt: time.Now(),
		CachedObj: value,
	}

}

func (cache *CacheService) CleanValue(key any) {

	cache.mu.Lock()
	defer cache.mu.Unlock()

	delete(cache.store, key)
}

func (cache *CacheService) EraseValues() {

	cache.mu.Lock()
	defer cache.mu.Unlock()

	for key := range cache.store {
		delete(cache.store, key)
	}
}

func (cache *CacheService) OlderThanAndExists(key any, duration time.Duration) bool {

	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, exists := cache.store[key]
	return (time.Since(val.CreatedAt) > duration) && exists
}

func (cache *CacheService) CleanOlderThan(duration time.Duration) {

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	for key, val := range cache.store {
		if time.Since(val.CreatedAt) > duration {
			delete(cache.store, key)
		}
	}
}
