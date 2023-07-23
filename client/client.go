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

type Config struct {
	AccessTokenClient AccessTokenClient // The client used for request access token
	BaseApiUrl        *url.URL          // The endpoint to request an API, if full path is not given, default value is 'https://api.weixin.qq.com'
	Cache             caches.Cache      // Cache instance for managing tokens
	HttpClient        *http.Client      // Default Http client to send request
}

// Represents an error that occurs when the WeChat API returns an unexpected code.
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
	GetAccessToken() (*Token, error)

	// FetchAccessToken renews and retrieves an access token.
	// It may return an error along with the token if there is no `Cache` set up.
	FetchAccessToken() (*Token, error)
}

type weChatClient struct {
	auth   auth.Auth
	cm     CredentialManager
	cache  caches.Cache
	client *http.Client
}

// Create a new `WeChatClient`
func New(auth auth.Auth, conf Config) WeChatClient {
	var (
		atc        AccessTokenClient
		baseApiUrl *url.URL
		client     http.Client
	)

	if conf.BaseApiUrl == nil {
		baseApiUrl, _ = url.Parse(DefaultBaseApiUri)
	} else {
		baseApiUrl = conf.BaseApiUrl
	}

	if conf.HttpClient == nil {
		client = http.Client{}
	} else {
		client = *conf.HttpClient
	}

	if conf.AccessTokenClient == nil {
		accessTokenUri, _ := url.Parse(DefaultAccessTokenUri)
		atcClient := client
		atc = NewAccessTokenClient(
			accessTokenUri,
			&authCredentialManager{auth: auth},
			&atcClient,
		)
	} else {
		atc = conf.AccessTokenClient
	}

	cm := &accessTokenCredentialManager{
		atc:   atc,
		auth:  auth,
		cache: conf.Cache,
	}
	client.Transport = NewCredentialRoundTripper(NewAccessTokenRoundTripper(NewCommonRoundTripper(client.Transport, baseApiUrl)), cm)

	return &weChatClient{
		cm:     cm,
		auth:   auth,
		cache:  conf.Cache,
		client: &client,
	}
}

func (c *weChatClient) Get(url string, withCredential bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, withCredential)
}

func (c *weChatClient) PostJson(url string, data interface{}, withCredential bool) (*http.Response, error) {
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

func (c *weChatClient) Do(req *http.Request, withCredential bool) (*http.Response, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestContextWithCredential, withCredential)
	req = req.WithContext(ctx)

	return c.client.Do(req)
}

func (c *weChatClient) GetAuth() auth.Auth {
	return c.auth
}

func (c *weChatClient) GetAccessToken() (*Token, error) {
	token, err := c.cm.Get()
	if token != nil {
		return token.(*Token), err
	}
	return nil, err
}

func (c *weChatClient) FetchAccessToken() (*Token, error) {
	token, err := c.cm.Renew()
	if token != nil {
		return token.(*Token), err
	}
	return nil, err
}
