package client

import "errors"

var (
	ErrCacheNotSet     = errors.New("cache not set")
	ErrInvalidResponse = errors.New("invalid response")
)
