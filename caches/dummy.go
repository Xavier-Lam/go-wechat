package caches

import (
	"time"
)

type cacheItem struct {
	Value     []byte
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

func (c *dummyCache) Get(appId string, key string) ([]byte, error) {
	fullKey := getKey(appId, key)
	item, ok := c.cache[fullKey]
	if !ok {
		return nil, ErrKeyNotFound
	}

	if item.ExpiresAt.Before(time.Now()) {
		delete(c.cache, fullKey)
		return nil, ErrKeyNotFound
	}

	return item.Value, nil
}

func (c *dummyCache) Set(appId string, key string, value []byte, expiresIn int) error {
	fullKey := getKey(appId, key)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	item := cacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	c.cache[fullKey] = item
	return nil
}

func (c *dummyCache) Add(appId string, key string, value []byte, expiresIn int) error {
	fullKey := getKey(appId, key)
	currentValue, err := c.Get(appId, key)
	if err != nil && err != ErrKeyNotFound {
		return err
	}
	if currentValue != nil {
		return ErrKeyExisted
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	item := cacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	c.cache[fullKey] = item
	return nil
}

func (c *dummyCache) Delete(appId string, key string, value []byte) error {
	fullKey := getKey(appId, key)
	if value != nil {
		currentValue, err := c.Get(appId, key)
		if err == ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		if string(currentValue) != string(value) {
			return ErrValueNotMatched
		}
	}
	delete(c.cache, fullKey)
	return nil
}

func getKey(appId, biz string) string {
	return appId + ":" + biz
}
