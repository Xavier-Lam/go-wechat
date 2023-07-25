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
	// AccessTokenClient is used for request a latest access token when it is needed
	// This option should be left as the default value (nil), unless you want to customize the client
	// For example, if you want to request your access token from a different service than Tencent's.
	AccessTokenClient auth.AccessTokenClient

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
	cm     auth.CredentialManager
	client *http.Client
}

// Create a new `WeChatClient`
func New(a auth.Auth, conf Config) WeChatClient {
	var (
		atc        auth.AccessTokenClient
		baseApiUrl *url.URL
		baseClient http.Client
	)

	if conf.BaseApiUrl == nil {
		baseApiUrl, _ = url.Parse(DefaultBaseApiUri)
	} else {
		baseApiUrl = conf.BaseApiUrl
	}

	if conf.HttpClient == nil {
		baseClient = http.Client{Transport: http.DefaultTransport}
	} else {
		baseClient = *conf.HttpClient
	}

	if conf.AccessTokenClient == nil {
		client := baseClient
		atc = NewAccessTokenClient(&client, "")
	} else {
		atc = conf.AccessTokenClient
	}

	cm := auth.NewAccessTokenManager(atc, a, conf.Cache)

	baseClient.Transport =
		NewCredentialRoundTripper(cm,
			NewAccessTokenRoundTripper(
				NewCommonRoundTripper(
					baseApiUrl, baseClient.Transport)))

	return &DefaultWeChatClient{
		cm:     cm,
		auth:   a,
		client: &baseClient,
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
	token, err := c.cm.Get()
	if token != nil {
		return token.(*auth.AccessToken), err
	}
	return nil, err
}

func (c *DefaultWeChatClient) FetchAccessToken() (*auth.AccessToken, error) {
	token, err := c.cm.Renew()
	if token != nil {
		return token.(*auth.AccessToken), err
	}
	return nil, err
}
