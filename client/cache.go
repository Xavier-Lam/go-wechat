package client

import (
	"errors"
)

const (
	CacheDefaultKeyPrefix = "wx:"

	CacheBizAccessToken = "ak"
	CacheBizJSTicket    = "js_ticket"
)

var CacheErrorKeyNotFound = errors.New("Key not found in cache")

type Cache interface {
	// Get retrieves the value associated with the given appId and biz from the cache.
	// If successful, it returns the value as an interface{} and nil error.
	// If no value is found or an error occurs, it returns nil value and `CacheErrorKeyNotFound`.
	Get(appId string, biz string) (interface{}, error)

	// Set stores the given value in the cache with the provided expiration time for the specified appId and biz.
	// If successful, it returns nil error.
	// If an error occurs during the storing process, it returns an error containing details of the failure.
	Set(appId string, biz string, value interface{}, expiresIn int) error
}
