package cache

import "time"

var instance *CacheProvider

type (
	CacheProvider struct {
		Cache ProvideMap
	}

	ProvideMap map[any]CacheDataUnit

	CacheDataUnit struct {
		CreatedAt time.Time
		CachedObj any

		_ struct{}
	}
)

func InitCache() {
	instance = &CacheProvider{
		Cache: ProvideMap{},
	}
}

func GetEntry() *CacheProvider {
	return instance
}

func (cache *CacheProvider) GetValue(key any) any {
	return cache.Cache[key].CachedObj
}

func (cache *CacheProvider) PushValue(key, value any) {
	cache.Cache[key] = CacheDataUnit{
		CreatedAt: time.Now(),
		CachedObj: value,
	}
}

func (cache *CacheProvider) OlderThanAndExists(key any, duration time.Duration) bool {
	val, exists := cache.Cache[key]
	return (time.Since(val.CreatedAt) > duration) && !exists
}
