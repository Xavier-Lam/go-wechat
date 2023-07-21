package caches_test

import (
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/stretchr/testify/assert"
)

type CacheFactory func() caches.Cache

func TestDummyCache(t *testing.T) {
	testCache(t, func() caches.Cache {
		return caches.NewDummyCache()
	})
}

// func TestRedisCache(t *testing.T) {
// 	testCache(t, func() w.Cache {
// 		s, err := miniredis.Run()
// 		assert.NoError(t, err)
// 		r := redis.NewClient(&redis.Options{
// 			Addr: s.Addr(),
// 		})
// 		return caches.NewRedisCache(r, "")
// 	})
// }

func testCache(t *testing.T, f CacheFactory) {
	testCacheGet(t, f)
	testCacheSet(t, f)
	testCacheAdd(t, f)
	testCacheDelete(t, f)
}

func testCacheGet(t *testing.T, f CacheFactory) {
	cache := f()

	// Set a value in the cache
	appId := "myAppId"
	key := "myBiz"
	value := []byte("myValue")
	expiresIn := 1 // seconds

	nilValue, err := cache.Get(appId, key)
	assert.ErrorIs(t, err, caches.ErrKeyNotFound)
	assert.Nil(t, nilValue, "Expected nil value")

	err = cache.Set(appId, key, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	// Get the value from the cache
	got, err := cache.Get(appId, key)
	assert.NoError(t, err, "Failed to get value from cache")

	// Check that the retrieved value matches the expected value
	assert.Equal(t, value, got, "Retrieved value does not match expected value")

	// Wait for the value to expire
	time.Sleep(time.Duration(expiresIn+1) * time.Second)

	// Get the expired value from the cache
	gotExpired, err := cache.Get(appId, key)
	assert.ErrorIs(t, err, caches.ErrKeyNotFound)
	assert.Nil(t, gotExpired, "Expected nil value")
}

func testCacheSet(t *testing.T, f CacheFactory) {
	cache := f()

	// Set a value in the cache
	appId := "myAppId"
	key := "myBiz"
	value := []byte("myValue")
	expiresIn := 10 // seconds
	err := cache.Set(appId, key, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	// Get the value from the cache
	got, err := cache.Get(appId, key)
	assert.NoError(t, err, "Failed to get value from cache")

	// Check that the retrieved value matches the expected value
	assert.Equal(t, value, got, "Retrieved value does not match expected value")
}

func testCacheAdd(t *testing.T, f CacheFactory) {
	cache := f()

	appId := "myAppId"
	key := "myBiz"
	value := []byte("myValue")
	value2 := []byte("myValue2")
	expiresIn := 1 // seconds
	err := cache.Add(appId, key, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	got, err := cache.Get(appId, key)
	assert.NoError(t, err, "Failed to get value from cache")
	assert.Equal(t, value, got, "Retrieved value does not match expected value")

	err = cache.Add(appId, key, value2, expiresIn)
	assert.ErrorIs(t, err, caches.ErrKeyExisted)

	got, err = cache.Get(appId, key)
	assert.NoError(t, err, "Failed to get value from cache")
	assert.Equal(t, value, got, "Retrieved value does not match expected value")

	// Wait for the value to expire
	time.Sleep(time.Duration(expiresIn+1) * time.Second)
	gotExpired, err := cache.Get(appId, key)
	assert.ErrorIs(t, err, caches.ErrKeyNotFound)
	assert.Nil(t, gotExpired, "Expected nil value")
}

func testCacheDelete(t *testing.T, f CacheFactory) {
	cache := f()

	appId := "myAppId"
	key := "myBiz"
	value := []byte("myValue")
	value2 := []byte("myValue2")
	expiresIn := 10 // seconds
	err := cache.Set(appId, key, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	// delete unmatched value
	err = cache.Delete(appId, key, value2)
	assert.ErrorIs(t, err, caches.ErrValueNotMatched)

	got, err := cache.Get(appId, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got, "Retrieved value does not match expected value")

	// delete matched value
	err = cache.Delete(appId, key, value)
	assert.NoError(t, err)

	got, err = cache.Get(appId, key)
	assert.ErrorIs(t, err, caches.ErrKeyNotFound)
	assert.Nil(t, got)

	// delete non-existed key
	err = cache.Delete(appId, key, value2)
	assert.NoError(t, err)
	err = cache.Delete(appId, key, nil)
	assert.NoError(t, err)

	// delete without value
	err = cache.Set(appId, key, value, expiresIn)
	assert.NoError(t, err)

	err = cache.Delete(appId, key, nil)
	assert.NoError(t, err)

	got, err = cache.Get(appId, key)
	assert.ErrorIs(t, err, caches.ErrKeyNotFound)
	assert.Nil(t, got)
}
