package client

import (
	"errors"
	"fmt"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

type CredentialManager interface {
	// Get the latest credential
	Get() (interface{}, error)

	// Renew credential
	Renew() (interface{}, error)

	// Delete a credential
	Delete() error
}

var ErrCacheNotSet = errors.New("cache not set")

type weChatAccessTokenCredentialManager struct {
	atc   AccessTokenClient
	auth  auth.Auth
	cache caches.Cache
}

func (cm *weChatAccessTokenCredentialManager) Get() (interface{}, error) {
	cachedValue, err := cm.get()
	if err == nil {
		return cachedValue, nil
	}

	return cm.Renew()
}

func (cm *weChatAccessTokenCredentialManager) Renew() (interface{}, error) {
	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.atc.GetAccessToken(cm.auth)
	if err != nil {
		return nil, err
	}

	if cm.cache == nil {
		err = fmt.Errorf("cache is not set")
	} else {
		serializedToken, err := serializeToken(token)
		if err != nil {
			return nil, err
		}
		err = cm.cache.Set(
			cm.auth.GetAppId(),
			caches.BizAccessToken,
			serializedToken,
			token.GetExpiresIn(),
		)
	}

	return token, err
}

func (cm *weChatAccessTokenCredentialManager) Delete() error {
	// TODO: prevent concurrent fetching
	token, err := cm.get()
	if err != nil {
		return err
	}
	serializedToken, err := serializeToken(token)
	if err != nil {
		return err
	}
	return cm.cache.Delete(
		cm.auth.GetAppId(),
		caches.BizAccessToken,
		serializedToken,
	)
}

func (m *weChatAccessTokenCredentialManager) get() (*Token, error) {
	if m.cache == nil {
		return nil, ErrCacheNotSet
	}

	cachedValue, err := m.cache.Get(m.auth.GetAppId(), caches.BizAccessToken)
	if err != nil {
		return nil, err
	}

	token, err := deserializeToken(cachedValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}
