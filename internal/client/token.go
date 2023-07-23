package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

const (
	DefaultAccessTokenUri = "https://api.weixin.qq.com/cgi-bin/token"
)

var ErrCacheNotSet = errors.New("cache not set")

type AccessTokenCredentialManager struct {
	atc   AccessTokenClient
	auth  auth.Auth
	cache caches.Cache
}

func (cm *AccessTokenCredentialManager) Get() (interface{}, error) {
	cachedValue, err := cm.get()
	if err == nil {
		return cachedValue, nil
	}

	return cm.Renew()
}

func (cm *AccessTokenCredentialManager) Set(credential interface{}) error {
	return errors.New("not settable")
}

func (cm *AccessTokenCredentialManager) Renew() (interface{}, error) {
	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.atc.GetAccessToken()
	if err != nil {
		return nil, err
	}

	if cm.cache == nil {
		err = fmt.Errorf("cache is not set")
	} else {
		serializedToken, err := auth.SerializeToken(token)
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

func (cm *AccessTokenCredentialManager) Delete() error {
	// TODO: prevent concurrent fetching
	token, err := cm.get()
	if err != nil {
		return err
	}
	serializedToken, err := auth.SerializeToken(token)
	if err != nil {
		return err
	}
	return cm.cache.Delete(
		cm.auth.GetAppId(),
		caches.BizAccessToken,
		serializedToken,
	)
}

func (m *AccessTokenCredentialManager) get() (*auth.AccessToken, error) {
	if m.cache == nil {
		return nil, ErrCacheNotSet
	}

	cachedValue, err := m.cache.Get(m.auth.GetAppId(), caches.BizAccessToken)
	if err != nil {
		return nil, err
	}

	token, err := auth.DeserializeToken(cachedValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}

type TokenResponse interface {
	GetAccessToken() string
	GetExpiresIn() int
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (t *tokenResponse) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenResponse) GetExpiresIn() int {
	return t.ExpiresIn
}

type AccessTokenClient interface {
	GetAccessToken() (*auth.AccessToken, error)
}

type TokenClient struct {
	client   *http.Client
	endpoint *url.URL // The endpoint to request a new token, default value is 'https://api.weixin.qq.com/cgi-bin/token'
	dto      TokenResponse
}

func NewAccessTokenClient(endpoint *url.URL, cm auth.CredentialManager, client *http.Client) AccessTokenClient {
	if client == nil {
		client = &http.Client{}
	}
	client.Transport =
		NewCredentialRoundTripper(cm,
			NewFetchAccessTokenRoundTripper(
				NewCommonRoundTripper(nil, client.Transport)))

	return &TokenClient{
		client:   client,
		endpoint: endpoint,
		dto:      &tokenResponse{},
	}
}

func (c *TokenClient) GetAccessToken() (*auth.AccessToken, error) {
	// Prepare request
	req, err := http.NewRequest(http.MethodGet, c.endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestContextWithCredential, true)
	req = req.WithContext(ctx)

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse token
	token := c.dto
	err = GetJson(resp, token)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if token.GetAccessToken() == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	rv := auth.NewAccessToken(token.GetAccessToken(), token.GetExpiresIn())
	return rv, nil
}
