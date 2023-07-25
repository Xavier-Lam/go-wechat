package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

const (
	DefaultAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type accessTokenClient struct {
	client     *http.Client
	requestUrl *url.URL
}

// NewAccessTokenClient creates an `auth.AccessTokenClient` to get the latest access token
func NewAccessTokenClient(baseClient *http.Client, rawRequestUrl string) auth.AccessTokenClient {
	if baseClient == nil {
		baseClient = &http.Client{Transport: http.DefaultTransport}
	}

	client := *baseClient
	client.Transport = NewCommonRoundTripper(nil, client.Transport)

	if rawRequestUrl == "" {
		rawRequestUrl = DefaultAccessTokenUrl
	}
	requestUrl, err := url.Parse(rawRequestUrl)
	if err != nil {
		requestUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	return &accessTokenClient{
		client:     &client,
		requestUrl: requestUrl,
	}
}

func (c *accessTokenClient) PrepareRequest(a auth.Auth) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, c.requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{
		"grant_type": {"client_credential"},
		"appid":      {a.GetAppId()},
		"secret":     {a.GetAppSecret()},
	}
	req.URL.RawQuery = query.Encode()

	return req, nil
}

func (c *accessTokenClient) SendRequest(a auth.Auth, req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *accessTokenClient) HandleResponse(a auth.Auth, resp *http.Response, req *http.Request) (*auth.AccessToken, error) {
	data := &tokenResponse{}
	err := GetJson(resp, data)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if data.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	token := auth.NewAccessToken(data.AccessToken, data.ExpiresIn)
	return token, nil
}
