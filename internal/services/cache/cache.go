package cache

import "time"

var instance *CacheService

type (
	// Implements Global value cache for web-server
	CacheService struct {
		Cache ProvideMap
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
		Cache: ProvideMap{},
	}
}

func GetEntry() *CacheService {
	return instance
}

func (cache *CacheService) GetValue(key any) any {
	return cache.Cache[key].CachedObj
}

func (cache *CacheService) PushValue(key, value any) {
	cache.Cache[key] = CacheDataUnit{
		CreatedAt: time.Now(),
		CachedObj: value,
	}
}

func (cache *CacheService) OlderThanAndExists(key any, duration time.Duration) bool {
	val, exists := cache.Cache[key]
	return (time.Since(val.CreatedAt) > duration) && exists
}
