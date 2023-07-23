package client

import (
	"errors"
	"fmt"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

var ErrCacheNotSet = errors.New("cache not set")

// It would be much better if Go supports covariance...
type CredentialManager interface {
	// Get the latest credential
	Get() (interface{}, error)

	// Set the latest credential
	Set(credential interface{}) error

	// Renew credential
	Renew() (interface{}, error)

	// Delete a credential
	Delete() error
}

type authCredentialManager struct {
	auth auth.Auth
}

func (cm *authCredentialManager) Get() (interface{}, error) {
	if cm.auth == nil {
		return errors.New("auth not set"), nil
	}
	return cm.auth, nil
}

func (cm *authCredentialManager) Set(credential interface{}) error {
	return errors.New("not settable")
}

func (cm *authCredentialManager) Renew() (interface{}, error) {
	return nil, errors.New("not renewable")
}

func (cm *authCredentialManager) Delete() error {
	return errors.New("not deletable")
}

type accessTokenCredentialManager struct {
	atc   AccessTokenClient
	auth  auth.Auth
	cache caches.Cache
}

func (cm *accessTokenCredentialManager) Get() (interface{}, error) {
	cachedValue, err := cm.get()
	if err == nil {
		return cachedValue, nil
	}

	return cm.Renew()
}

func (cm *accessTokenCredentialManager) Set(credential interface{}) error {
	return errors.New("not settable")
}

func (cm *accessTokenCredentialManager) Renew() (interface{}, error) {
	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.atc.GetAccessToken()
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

func (cm *accessTokenCredentialManager) Delete() error {
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

func (m *accessTokenCredentialManager) get() (*Token, error) {
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
