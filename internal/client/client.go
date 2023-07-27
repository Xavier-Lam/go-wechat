package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

const (
	DefaultBaseApiUri = "https://api.weixin.qq.com"

	ErrCodeInvalidCredential  = 40001
	ErrCodeInvalidAccessToken = 40014
	ErrCodeAccessTokenExpired = 42001

	RequestContextWithCredential = "withCredential"
	RequestContextCredential     = "credential"
)

// WeChatClient is an interface representing a client for making requests to WeChat APIs.
type WeChatClient interface {
	// Sends a GET request
	Get(url string, withCredential bool) (*http.Response, error)

	// Sends a POST request with JSON data
	PostJson(url string, data interface{}, withCredential bool) (*http.Response, error)

	// Sends a request with or without access_token based on `withCredential` flag
	Do(req *http.Request, withCredential bool) (*http.Response, error)

	// Returns the WeChat Auth of the client
	GetAuth() auth.Auth

	// GetAccessToken retrieves the access token.
	// It may return an error along with the token if there is no `Cache` set up.
	GetAccessToken() (*auth.AccessToken, error)

	// FetchAccessToken renews and retrieves an access token.
	// It may return an error along with the token if there is no `Cache` set up.
	FetchAccessToken() (*auth.AccessToken, error)
}

// Config is a configuration struct used to set up a `client.WeChatClient`.
type Config struct {
	// AccessTokenFetcher is a callback function to return the latest access token
	// The default implement should be suitable for most case, override only when
	// you want to customize the way you make request.
	// For example, if you want to request to a service rather than Tencent's.
	AccessTokenFetcher AccessTokenFetcher

	// AccessTokenUrl is the url AccessTokenFetcher tries to fetch the latest access token.
	// This URL will be passed to the AccessTokenFetcher callback.
	// If not provided, the default value is 'https://api.weixin.qq.com/cgi-bin/token'.
	AccessTokenUrl *url.URL

	// BaseApiUrl is the base URL used for making API requests.
	// If not provided, the default value is 'https://api.weixin.qq.com'.
	BaseApiUrl *url.URL

	// Cache instance for managing tokens
	Cache caches.Cache

	// HttpClient is the default HTTP client used for sending requests.
	HttpClient *http.Client
}

// WeChatApiError represents an error that occurs when the WeChat API returns an unexpected code.
type WeChatApiError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// Error happened while retrying
	RetryError error
}

func (e WeChatApiError) Error() string {
	if e.RetryError != nil {
		return fmt.Sprintf("WeChat API error [%d]: %s (Retry error: %s)", e.ErrCode, e.ErrMsg, e.RetryError.Error())
	} else {
		return fmt.Sprintf("WeChat API error [%d]: %s", e.ErrCode, e.ErrMsg)
	}
}

type DefaultWeChatClient struct {
	auth   auth.Auth
	cm     auth.AccessTokenManager
	client *http.Client
}

// Create a new `WeChatClient`
func New(a auth.Auth, conf Config) WeChatClient {
	if conf.HttpClient == nil {
		conf.HttpClient = &http.Client{Transport: http.DefaultTransport}
	}

	if conf.AccessTokenFetcher == nil {
		conf.AccessTokenFetcher = accessTokenFetcher
	}

	if conf.AccessTokenUrl == nil {
		conf.AccessTokenUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	fetcher := func() (*auth.AccessToken, error) {
		c := *conf.HttpClient
		return conf.AccessTokenFetcher(&c, a, conf.AccessTokenUrl)
	}
	cm := auth.NewAccessTokenManager(a, conf.Cache, fetcher)

	if conf.BaseApiUrl == nil {
		conf.BaseApiUrl, _ = url.Parse(DefaultBaseApiUri)
	}

	c := *conf.HttpClient
	c.Transport =
		NewCredentialRoundTripper(cm,
			NewAccessTokenRoundTripper(
				NewCommonRoundTripper(
					conf.BaseApiUrl, conf.HttpClient.Transport)))

	return &DefaultWeChatClient{
		cm:     cm,
		auth:   a,
		client: &c,
	}
}

func (c *DefaultWeChatClient) Get(url string, withCredential bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, withCredential)
}

func (c *DefaultWeChatClient) PostJson(url string, data interface{}, withCredential bool) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.Do(req, withCredential)
}

func (c *DefaultWeChatClient) Do(req *http.Request, withCredential bool) (*http.Response, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestContextWithCredential, withCredential)
	req = req.WithContext(ctx)

	return c.client.Do(req)
}

func (c *DefaultWeChatClient) GetAuth() auth.Auth {
	return c.auth
}

func (c *DefaultWeChatClient) GetAccessToken() (*auth.AccessToken, error) {
	return c.cm.Get()
}

func (c *DefaultWeChatClient) FetchAccessToken() (*auth.AccessToken, error) {
	return c.cm.Renew()
}
