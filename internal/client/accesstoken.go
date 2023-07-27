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

type AccessTokenFetcher = func(client *http.Client, auth auth.Auth, accessTokenUrl *url.URL) (*auth.AccessToken, error)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// accessTokenFetcher is a callback function used for getting the latest access token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func accessTokenFetcher(c *http.Client, a auth.Auth, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
	// Validate parameters
	appId := a.GetAppId()
	appSecret := a.GetAppSecret()
	if appId == "" || appSecret == "" {
		return nil, fmt.Errorf("invalid auth: %s", appId)
	}

	if accessTokenUrl == nil {
		accessTokenUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	// Set up HTTP client
	if c == nil {
		c = &http.Client{Transport: http.DefaultTransport}
	}
	c.Transport = NewCommonRoundTripper(nil, c.Transport)

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, accessTokenUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{
		"grant_type": {"client_credential"},
		"appid":      {appId},
		"secret":     {appSecret},
	}
	req.URL.RawQuery = query.Encode()

	// Send request
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	data := &tokenResponse{}
	err = GetJson(resp, data)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if data.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	token := auth.NewAccessToken(data.AccessToken, data.ExpiresIn)
	return token, nil
}
