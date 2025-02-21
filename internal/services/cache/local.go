package cache

type LocalCache struct {
	CacheService
}

func NewLocalCache() *LocalCache {
	return &LocalCache{
		CacheService{
			store: ProvideMap{},
		},
	}
}
