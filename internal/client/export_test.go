package client

import (
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

func NewWeChatAccessTokenCredentialManager(auth auth.Auth, cache caches.Cache, akc AccessTokenClient) auth.CredentialManager {
	return &AccessTokenCredentialManager{
		atc:   akc,
		auth:  auth,
		cache: cache,
	}
}
