package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

const (
	DefaultAccessTokenUri = "https://api.weixin.qq.com/cgi-bin/token"
	DefaultTokenExpiresIn = 7200
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Token struct {
	accessToken string
	expiresIn   int
	createdAt   time.Time
}

func NewToken(accessToken string, expiresIn int) *Token {
	if expiresIn <= 0 {
		expiresIn = DefaultTokenExpiresIn
	}
	return &Token{
		accessToken: accessToken,
		expiresIn:   expiresIn,
		createdAt:   time.Now(),
	}
}

func (t *Token) GetAccessToken() string {
	return t.accessToken
}

func (t *Token) GetExpiresIn() int {
	timeDiff := time.Since(t.createdAt)
	timeEscaped := int(timeDiff.Seconds())
	if timeEscaped >= t.expiresIn {
		return 0
	}
	return t.expiresIn - timeEscaped
}

func (t *Token) GetExpiresAt() time.Time {
	timeDiff := time.Duration(t.expiresIn) * time.Second
	return t.createdAt.Add(timeDiff)
}

type AccessTokenClient interface {
	GetAccessToken(auth auth.Auth) (*Token, error)
}

type accessTokenClient struct {
	client   *http.Client
	endpoint *url.URL // The endpoint to request a new token, default value is 'https://api.weixin.qq.com/cgi-bin/token'
}

func NewAccessTokenClient(endpoint *url.URL, client *http.Client) AccessTokenClient {
	if client == nil {
		client = &http.Client{
			Transport: NewCommonRoundTripper(http.DefaultTransport, nil),
		}
	}
	return &accessTokenClient{
		client:   client,
		endpoint: endpoint,
	}
}

func (c *accessTokenClient) GetAccessToken(auth auth.Auth) (*Token, error) {
	// Build url
	uri := c.endpoint.String()
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
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse token
	token := &tokenResponse{}
	err = GetJson(resp, &token)
	if err != nil {
		return nil, fmt.Errorf("malformed access token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token response")
	}

	rv := NewToken(token.AccessToken, token.ExpiresIn)
	return rv, nil
}

type serializedTokenData struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	CreatedAt   time.Time `json:"created_at"`
}

func serializeToken(token *Token) ([]byte, error) {
	timeDiff := -time.Duration(time.Second * time.Duration(token.GetExpiresIn()))
	data := &serializedTokenData{
		AccessToken: token.GetAccessToken(),
		ExpiresIn:   token.GetExpiresIn(),
		CreatedAt:   token.GetExpiresAt().Add(timeDiff),
	}
	return json.Marshal(data)
}

func deserializeToken(bytes []byte) (*Token, error) {
	data := &serializedTokenData{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &Token{
		accessToken: data.AccessToken,
		expiresIn:   data.ExpiresIn,
		createdAt:   data.CreatedAt,
	}, nil
}
