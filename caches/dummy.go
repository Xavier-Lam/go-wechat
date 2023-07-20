package caches

import (
	"time"
)

type cacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type dummyCache struct {
	cache map[string]cacheItem
}

func NewDummyCache() Cache {
	return &dummyCache{
		cache: make(map[string]cacheItem),
	}
}

func (c *dummyCache) Get(appId string, biz string) (interface{}, error) {
	key := getKey(appId, biz)
	item, ok := c.cache[key]
	if !ok {
		return nil, CacheErrorKeyNotFound
	}

	if item.ExpiresAt.Before(time.Now()) {
		delete(c.cache, key)
		return nil, CacheErrorKeyNotFound
	}

	return item.Value, nil
}

func (c *dummyCache) Set(appId string, biz string, value interface{}, expiresIn int) error {
	key := getKey(appId, biz)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	item := cacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	c.cache[key] = item
	return nil
}

func getKey(appId, biz string) string {
	return appId + ":" + biz
}
