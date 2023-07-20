package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Xavier-Lam/go-wechat"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	HttpClient     HttpClient // Default Http client to send request
	Cache          Cache      // Cache instance for managing tokens
	AccessTokenUri *url.URL   // The endpoint to request a new token, default value is 'https://api.weixin.qq.com/cgi-bin/token'
	BaseApiUri     *url.URL   // The endpoint to request an API, if full path is not given, default value is 'https://api.weixin.qq.com'
}

const (
	DefaultBaseApiUri = "https://api.weixin.qq.com"

	ErrCodeInvalidCredential  = 40001
	ErrCodeInvalidAccessToken = 40014
	ErrCodeAccessTokenExpired = 42001
)

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
	GetAuth() wechat.Auth

	// GetAccessToken retrieves the access token.
	// It may return an error along with the token if there is no `Cache` set up.
	GetAccessToken() (Token, error)

	// FetchAccessToken renews and retrieves an access token.
	// It may return an error along with the token if there is no `Cache` set up.
	FetchAccessToken() (Token, error)
}

type weChatClient struct {
	akc     AccessTokenClient
	auth    wechat.Auth
	baseUri *url.URL
	cache   Cache
	http    HttpClient
}

// Create a new `WeChatClient`
func New(auth wechat.Auth, conf *Config) WeChatClient {
	// Set up base dependencies if not given
	if conf == nil {
		conf = &Config{}
	}
	if conf.AccessTokenUri == nil {
		conf.AccessTokenUri, _ = url.Parse(DefaultAccessTokenUri)
	}
	if conf.HttpClient == nil {
		conf.HttpClient = &http.Client{}
	}
	akc := &accessTokenClient{
		http: conf.HttpClient,
		url:  conf.AccessTokenUri,
	}

	return NewWithDependency(akc, auth, conf)
}

// If you want to inject your own dependencies, you can use this constructor to create a client
func NewWithDependency(akc AccessTokenClient, auth wechat.Auth, conf *Config) WeChatClient {
	// Set up base dependencies if not given
	if conf == nil {
		conf = &Config{}
	}
	if conf.BaseApiUri == nil {
		conf.BaseApiUri, _ = url.Parse(DefaultBaseApiUri)
	}
	if conf.HttpClient == nil {
		conf.HttpClient = &http.Client{}
	}

	return &weChatClient{
		akc:     akc,
		auth:    auth,
		baseUri: conf.BaseApiUri,
		cache:   conf.Cache,
		http:    conf.HttpClient,
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
	ctx := req.Context()
	ctx = context.WithValue(ctx, "token", nil)

	// Prepare request
	if !req.URL.IsAbs() {
		req.URL = c.baseUri.ResolveReference(req.URL)
	}

	if withCredential {
		token, err := c.GetAccessToken()
		if token == nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "token", token)
	}

	req = req.WithContext(ctx)

	// Send request
	resp, err := c.do(req)
	if err == nil {
		return resp, nil
	}

	// Handle exception
	return c.handleError(err, req, resp)
}

func (c *weChatClient) GetAuth() wechat.Auth {
	return c.auth
}

func (c *weChatClient) GetAccessToken() (Token, error) {
	if c.cache != nil {
		cachedValue, err := c.cache.Get(c.auth.GetAppId(), "ak")
		if err == nil {
			if token, ok := cachedValue.(Token); ok {
				return token, nil
			}
		}
	}

	return c.FetchAccessToken()
}

func (c *weChatClient) FetchAccessToken() (Token, error) {
	token, err := c.akc.GetAccessToken(c.auth)
	if err != nil {
		return nil, err
	}

	if c.cache == nil {
		err = fmt.Errorf("Cache is not set")
	} else {
		err = c.cache.Set(c.auth.GetAppId(), "ak", token, token.GetExpiresIn())
	}

	return token, err
}

func (c *weChatClient) do(req *http.Request) (*http.Response, error) {
	token := req.Context().Value("token")
	if token != nil {
		query := req.URL.Query()
		query.Set("access_token", token.(Token).GetAccessToken())
		req.URL.RawQuery = query.Encode()
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error occurred when sending request: %w", err)
	}

	err = processResponse(resp)
	return resp, err
}

func (c *weChatClient) handleError(err error, req *http.Request, resp *http.Response) (*http.Response, error) {
	defer resp.Body.Close()

	apiError, ok := err.(WeChatApiError)
	if !ok {
		return nil, err
	}

	switch apiError.ErrCode {
	case ErrCodeAccessTokenExpired,
		ErrCodeInvalidAccessToken,
		ErrCodeInvalidCredential:
		ctx := req.Context()
		if ctx.Value("token") != nil {
			var token Token
			token, apiError.RetryError = c.FetchAccessToken()
			if token == nil {
				return nil, apiError
			}

			ctx = context.WithValue(ctx, "token", token)
			req = req.WithContext(ctx)
			return c.do(req)
		}
	}

	return nil, apiError
}

func processResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var apiError WeChatApiError
		err := GetJson(resp, &apiError)
		if err != nil {
			return fmt.Errorf("Failed to decode JSON: %w", err)
		} else if apiError.ErrCode != 0 {
			return apiError
		}
	}

	return nil
}

func GetJson(resp *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return err
}
