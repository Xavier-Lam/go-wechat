package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Xavier-Lam/go-wechat"
)

const (
	DefaultAccessTokenUri = "https://api.weixin.qq.com/cgi-bin/token"
	DefaultTokenExpiresIn = 7200
)

type TokenData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Token interface {
	GetAccessToken() string
	GetExpiresIn() int
	GetExpiresAt() time.Time
}

type token struct {
	accessToken string
	expiresIn   int
	createdAt   time.Time
}

func NewToken(accessToken string, expiresIn int) Token {
	if expiresIn <= 0 {
		expiresIn = DefaultTokenExpiresIn
	}
	return &token{
		accessToken: accessToken,
		expiresIn:   expiresIn,
		createdAt:   time.Now(),
	}
}

func (t *token) GetAccessToken() string {
	return t.accessToken
}

func (t *token) GetExpiresIn() int {
	timeDiff := time.Since(t.createdAt)
	timeEscaped := int(timeDiff.Seconds())
	if timeEscaped >= t.expiresIn {
		return 0
	}
	return t.expiresIn - timeEscaped
}

func (t *token) GetExpiresAt() time.Time {
	timeDiff := time.Duration(t.expiresIn) * time.Second
	return t.createdAt.Add(timeDiff)
}

type AccessTokenClient interface {
	GetAccessToken(auth wechat.Auth) (Token, error)
}

type accessTokenClient struct {
	http HttpClient
	url  *url.URL
}

func (c *accessTokenClient) GetAccessToken(auth wechat.Auth) (Token, error) {
	// Build url
	uri := c.url.String()
	query := url.Values{
		"grant_type": {"client_credential"},
		"appid":      {auth.GetAppId()},
		"secret":     {auth.GetAppSecret()},
	}
	uri += "?" + query.Encode()

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	// Send request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error occurred when sending request: %w", err)
	}

	// Handle response
	defer resp.Body.Close()
	err = processResponse(resp)
	if err != nil {
		return nil, err
	}

	// Parse token
	token := &TokenData{}
	err = GetJson(resp, &token)
	if err != nil {
		return nil, fmt.Errorf("Malformed access token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("Failed to retrieve access token: Invalid response")
	}

	rv := NewToken(token.AccessToken, token.ExpiresIn)
	return rv, nil
}
