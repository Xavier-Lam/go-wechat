package caches

import "errors"

const (
	DefaultKeyPrefix = "wx:"

	BizAccessToken = "ak"
	BizJSTicket    = "js_ticket"
)

var (
	ErrKeyExisted      = errors.New("key existed in cache")
	ErrKeyNotFound     = errors.New("key not found in cache")
	ErrValueNotMatched = errors.New("value not matched")
	ErrCacheNotSet     = errors.New("cache not set")
)

type Cache interface {
	// Get retrieves the value associated with the given appId and biz from the cache.
	// If successful, it returns the value as an interface{} and nil error.
	// If no value is found or an error occurs, it returns nil value and `CacheErrorKeyNotFound`.
	Get(appId string, key string) ([]byte, error)

	// Set stores the given value in the cache with the provided expiration time for the specified appId and biz.
	// If successful, it returns nil error.
	// If an error occurs during the storing process, it returns an error containing details of the failure.
	Set(appId string, key string, value []byte, expiresIn int) error

	// Add value to cache if key not existed
	Add(appId string, key string, value []byte, expiresIn int) error

	// Delete stored value given
	// Compare the value before deleting if given
	Delete(appId string, key string, value []byte) error
}
