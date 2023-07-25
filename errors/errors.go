package errors

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
)

// Exported errors
type (
	WeChatApiError = client.WeChatApiError
)

var (
	ErrCacheNotSet     = auth.ErrCacheNotSet
	ErrInvalidResponse = client.ErrInvalidResponse
)
