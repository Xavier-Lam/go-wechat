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

// AccessTokenClient is an client to request the newest access token
type AccessTokenClient interface {
	GetAccessToken() (*auth.AccessToken, error)
}

// AccessTokenResponse represents the response data received from the server
// for an access token request.
type AccessTokenResponse interface {
	GetAccessToken() string
	GetExpiresIn() int
}

// AccessTokenManagerProvider is a factory function to create a `auth.CredentialManager`
// for manage access token of a WeChat application
type AccessTokenManagerProvider = func(
	auth auth.Auth,
	client http.Client,
	cache caches.Cache,
	accessTokenUrl *url.URL,
) auth.CredentialManager

// AccessTokenManager is an implement of the `auth.CredentialManager`
// which is used to manage access token credentials.
type AccessTokenManager struct {
	atc   AccessTokenClient
	auth  auth.Auth
	cache caches.Cache
}

// NewAccessTokenManager creates a new instance of `auth.CredentialManager`
// to manage access token credentials.
func NewAccessTokenManager(atc AccessTokenClient, auth auth.Auth, cache caches.Cache) auth.CredentialManager {
	return &AccessTokenManager{
		atc:   atc,
		auth:  auth,
		cache: cache,
	}
}

func (cm *AccessTokenManager) Get() (interface{}, error) {
	cachedValue, err := cm.get()
	if err == nil {
		return cachedValue, nil
	}

	return cm.Renew()
}

func (cm *AccessTokenManager) Set(credential interface{}) error {
	return errors.New("not settable")
}

func (cm *AccessTokenManager) Renew() (interface{}, error) {
	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.atc.GetAccessToken()
	if err != nil {
		return nil, err
	}

	if cm.cache == nil {
		err = fmt.Errorf("cache is not set")
	} else {
		serializedToken, err := auth.SerializeAccessToken(token)
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

func (cm *AccessTokenManager) Delete() error {
	token, err := cm.get()
	if err != nil {
		return err
	}
	serializedToken, err := auth.SerializeAccessToken(token)
	if err != nil {
		return err
	}
	return cm.cache.Delete(
		cm.auth.GetAppId(),
		caches.BizAccessToken,
		serializedToken,
	)
}

func (m *AccessTokenManager) get() (*auth.AccessToken, error) {
	if m.cache == nil {
		return nil, ErrCacheNotSet
	}

	cachedValue, err := m.cache.Get(m.auth.GetAppId(), caches.BizAccessToken)
	if err != nil {
		return nil, err
	}

	token, err := auth.DeserializeAccessToken(cachedValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func AccessTokenManagerFactory(auth auth.Auth, client http.Client, cache caches.Cache, accessTokenUrl *url.URL) auth.CredentialManager {
	atc := AccessTokenClientFactory(accessTokenUrl, auth, &client)
	return NewAccessTokenManager(atc, auth, cache)
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

type accessTokenClient struct {
	client     *http.Client
	requestUrl *url.URL // The url to request a new token, default value is 'https://api.weixin.qq.com/cgi-bin/token'
	dto        AccessTokenResponse
}

// NewAccessTokenClient creates the default access token client which is used to
// request the latest access token from server
func NewAccessTokenClient(client *http.Client, dto AccessTokenResponse, requestUrl *url.URL) AccessTokenClient {
	return &accessTokenClient{
		client:     client,
		requestUrl: requestUrl,
		dto:        &tokenResponse{},
	}
}

func (c *accessTokenClient) GetAccessToken() (*auth.AccessToken, error) {
	// Prepare request
	req, err := http.NewRequest(http.MethodGet, c.requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestContextWithCredential, true)
	req = req.WithContext(ctx)

	// Send request
	// Use `fetchAccessTokenRoundTripper` to set up request parameters
	// TODO: to use a request maker instead of RoundTrippers to send request (too complicated)
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

// AccessTokenClientFactory is a factory to creates the default access token client
// to request the latest access token from server
func AccessTokenClientFactory(requestUrl *url.URL, a auth.Auth, client *http.Client) AccessTokenClient {
	if client == nil {
		client = &http.Client{Transport: http.DefaultTransport}
	}
	client.Transport =
		NewCredentialRoundTripper(auth.NewAuthCredentialManager(a),
			NewFetchAccessTokenRoundTripper(
				NewCommonRoundTripper(nil, client.Transport)))

	if requestUrl == nil {
		requestUrl, _ = url.Parse(DefaultAccessTokenUri)
	}

	return NewAccessTokenClient(client, &tokenResponse{}, requestUrl)
}
