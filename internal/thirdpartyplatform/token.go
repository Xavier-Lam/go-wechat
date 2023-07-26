package thirdpartyplatform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
)

const (
	DefaultAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
)

type tokenResponse struct {
	AccessToken string `json:"component_access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type accessTokenClient struct {
	client     *http.Client
	requestUrl *url.URL
}

// NewAccessTokenClient creates an `auth.AccessTokenClient` to get the latest access token
// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
func NewAccessTokenClient(baseClient *http.Client, rawRequestUrl string) auth.AccessTokenClient {
	if baseClient == nil {
		baseClient = &http.Client{Transport: http.DefaultTransport}
	}

	c := *baseClient
	c.Transport = client.NewCommonRoundTripper(nil, c.Transport)

	if rawRequestUrl == "" {
		rawRequestUrl = DefaultAccessTokenUrl
	}
	requestUrl, err := url.Parse(rawRequestUrl)
	if err != nil {
		requestUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	return &accessTokenClient{
		client:     &c,
		requestUrl: requestUrl,
	}
}

func (c *accessTokenClient) PrepareRequest(a auth.Auth) (*http.Request, error) {
	tpa, ok := a.(ThirdPartyPlatformAuth)
	if !ok {
		return nil, errors.New("incorrect auth")
	}

	req, err := http.NewRequest(http.MethodPost, c.requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	ticket, err := tpa.GetTicket()
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"component_appid":         a.GetAppId(),
		"component_appsecret":     a.GetAppSecret(),
		"component_verify_ticket": ticket,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(http.MethodPost, req.URL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *accessTokenClient) SendRequest(a auth.Auth, req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *accessTokenClient) HandleResponse(a auth.Auth, resp *http.Response, req *http.Request) (*auth.AccessToken, error) {
	data := &tokenResponse{}
	err := client.GetJson(resp, data)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if data.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	token := auth.NewAccessToken(data.AccessToken, data.ExpiresIn)
	return token, nil
}
