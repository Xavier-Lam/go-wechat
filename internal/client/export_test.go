package client

import (
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

func NewWeChatAccessTokenCredentialManager(auth auth.Auth, cache caches.Cache, atc AccessTokenClient) auth.CredentialManager {
	return &AccessTokenCredentialManager{
		atc:   atc,
		auth:  auth,
		cache: cache,
	}
}
