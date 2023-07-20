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
}

func testCacheGet(t *testing.T, f CacheFactory) {
	cache := f()

	// Set a value in the cache
	appId := "myAppId"
	biz := "myBiz"
	value := "myValue"
	expiresIn := 1 // seconds

	nilValue, err := cache.Get(appId, biz)
	assert.ErrorIs(t, err, caches.CacheErrorKeyNotFound)
	assert.Nil(t, nilValue, "Expected nil value")

	err = cache.Set(appId, biz, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	// Get the value from the cache
	got, err := cache.Get(appId, biz)
	assert.NoError(t, err, "Failed to get value from cache")

	// Check that the retrieved value matches the expected value
	assert.Equal(t, value, got, "Retrieved value does not match expected value")

	// Wait for the value to expire
	time.Sleep(time.Duration(expiresIn+1) * time.Second)

	// Get the expired value from the cache
	gotExpired, err := cache.Get(appId, biz)
	assert.ErrorIs(t, err, caches.CacheErrorKeyNotFound)
	assert.Nil(t, gotExpired, "Expected nil value")
}

func testCacheSet(t *testing.T, f CacheFactory) {
	cache := f()

	// Set a value in the cache
	appId := "myAppId"
	biz := "myBiz"
	value := "myValue"
	expiresIn := 10 // seconds
	err := cache.Set(appId, biz, value, expiresIn)
	assert.NoError(t, err, "Failed to set value in cache")

	// Get the value from the cache
	got, err := cache.Get(appId, biz)
	assert.NoError(t, err, "Failed to get value from cache")

	// Check that the retrieved value matches the expected value
	assert.Equal(t, value, got, "Retrieved value does not match expected value")
}
