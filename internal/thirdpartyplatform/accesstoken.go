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

type AccessTokenFetcher = func(client *http.Client, auth auth.Auth, ticket string, accessTokenUrl *url.URL) (*auth.AccessToken, error)

type tokenResponse struct {
	AccessToken string `json:"component_access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// accessTokenFetcher is a callback function used for getting the latest access token
// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
func accessTokenFetcher(c *http.Client, a auth.Auth, ticket string, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
	// Validate parameters
	appId := a.GetAppId()
	appSecret := a.GetAppSecret()
	if appId == "" || appSecret == "" {
		return nil, fmt.Errorf("invalid auth: %s", appId)
	}

	if ticket == "" {
		return nil, errors.New("invalid ticket")
	}

	if accessTokenUrl == nil {
		accessTokenUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	// Set up HTTP client
	if c == nil {
		c = &http.Client{Transport: http.DefaultTransport}
	}
	c.Transport = client.NewCommonRoundTripper(nil, c.Transport)

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, accessTokenUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	reqData := map[string]string{
		"component_appid":         appId,
		"component_appsecret":     appSecret,
		"component_verify_ticket": ticket,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(http.MethodPost, req.URL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	data := &tokenResponse{}
	err = client.GetJson(resp, data)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if data.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	token := auth.NewAccessToken(data.AccessToken, data.ExpiresIn)
	return token, nil
}
