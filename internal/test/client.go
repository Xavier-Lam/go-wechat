package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/stretchr/testify/assert"
)

type RequestHandler func(req *http.Request, calls int) (*http.Response, error)

type mockHttpClient struct {
	calls   int
	handler RequestHandler
}

func NewMockHttpClient(handler RequestHandler) *http.Client {
	// return &mockHttpClient{handler: handler}
	return &http.Client{
		Transport: &mockHttpClient{handler: handler},
	}
}

func (c *mockHttpClient) RoundTrip(req *http.Request) (*http.Response, error) {
	c.calls++
	return c.handler(req, c.calls)
}

type responses struct{}

func (r *responses) Empty() (*http.Response, error) {
	return httptest.NewRecorder().Result(), nil
}

func (r *responses) Json(json string) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	recorder.Header().Add("Content-Type", "application/json")
	recorder.WriteString(json)
	return recorder.Result(), nil
}

var Responses = &responses{}

func AssertEndpointEqual(t *testing.T, expected string, actual *url.URL) {
	uri, err := url.Parse(expected)
	assert.NoError(t, err)
	assert.Equal(t, uri.Scheme, actual.Scheme)
	assert.Equal(t, uri.Host, actual.Host)
	assert.Equal(t, uri.Path, actual.Path)
}

type mockAccessTokenClient struct {
	token string
}

func NewMockAccessTokenClient(token string) auth.AccessTokenClient {
	return &mockAccessTokenClient{token: token}
}

func (c *mockAccessTokenClient) PrepareRequest(a auth.Auth) (*http.Request, error) {
	return &http.Request{}, nil
}

func (c *mockAccessTokenClient) SendRequest(a auth.Auth, req *http.Request) (*http.Response, error) {
	return Responses.Empty()
}

func (c *mockAccessTokenClient) HandleResponse(a auth.Auth, resp *http.Response, req *http.Request) (*auth.AccessToken, error) {
	return auth.NewAccessToken(c.token, auth.DefaultAccessTokenExpiresIn), nil
}
