package errors

import "github.com/Xavier-Lam/go-wechat/internal/client"

// Exported errors
type (
	WeChatApiError = client.WeChatApiError
)

var (
	ErrCacheNotSet = client.ErrCacheNotSet
)
