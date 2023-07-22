package client

import (
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

var (
	SerializeToken   = serializeToken
	DeserializeToken = deserializeToken
)

func NewWeChatAccessTokenCredentialManager(auth auth.Auth, cache caches.Cache, akc AccessTokenClient) CredentialManager {
	return &weChatAccessTokenCredentialManager{
		atc:   akc,
		auth:  auth,
		cache: cache,
	}
}
