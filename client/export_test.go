package client

import (
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

var (
	SerializeToken   = serializeToken
	DeserializeToken = deserializeToken
)

func NewAuthCredentialManager(auth auth.Auth) CredentialManager {
	return &authCredentialManager{auth: auth}
}

func NewWeChatAccessTokenCredentialManager(auth auth.Auth, cache caches.Cache, akc AccessTokenClient) CredentialManager {
	return &accessTokenCredentialManager{
		atc:   akc,
		auth:  auth,
		cache: cache,
	}
}
