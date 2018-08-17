package system

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheResolver func(key string) interface{}

type MemCache struct {
	cacheInstance *cache.Cache
	resolver CacheResolver
}

func (mc *MemCache) GetStore() (*cache.Cache) {
	return mc.cacheInstance
}

func (mc *MemCache) Remove(key string) {
	mc.cacheInstance.Delete(key)
}

func (mc *MemCache) SetResolver(resovler CacheResolver) {
	mc.resolver = resovler
}

func (mc *MemCache) Set(key string, value interface{}) {
	mc.cacheInstance.SetDefault(key, value)
}

func (mc *MemCache) Get(key string) (interface{}, bool){
	if mc.resolver == nil {
		return mc.TryGet(key)
	} else {
		return mc.EnsureGet(key)
	}
}

func (mc *MemCache) TryGet(key string) (interface{}, bool){
	return mc.cacheInstance.Get(key);
}

func (mc *MemCache) EnsureGet(key string) (interface{}, bool){
	value, ok := mc.cacheInstance.Get(key);
	if !ok && mc.resolver != nil {
		value = mc.resolver(key)
		ok = value != nil
		if (ok) {
			mc.Set(key, value)
		}
	}
	return value, ok
}

func CreateCache(defaultExpiration, cleanupInterval time.Duration) *MemCache {
	return &MemCache{cacheInstance: cache.New(defaultExpiration, cleanupInterval)}
}

func CreateEnsureCache(defaultExpiration, cleanupInterval time.Duration, resolver CacheResolver) *MemCache {
	return &MemCache{cacheInstance: cache.New(defaultExpiration, cleanupInterval), resolver: resolver}
}


