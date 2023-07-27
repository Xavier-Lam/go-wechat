package client_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenFetcher(t *testing.T) {
	httpClient := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, client.DefaultAccessTokenUrl, req.URL)
		assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
		assert.Equal(t, appID, req.URL.Query().Get("appid"))
		assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

		return test.Responses.Json(`{"access_token": "access-token", "expires_in": 7200}`)
	})

	token, err := client.ExportedAccessTokenFetcher(httpClient, test.MockAuth, nil)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())

	requestUrl, _ := url.Parse("https://example.com/test")
	httpClient = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, "https://example.com/test", req.URL)
		assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
		assert.Equal(t, appID, req.URL.Query().Get("appid"))
		assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

		return test.Responses.Json(`{"access_token": "access-token", "expires_in": 7200}`)
	})

	token, err = client.ExportedAccessTokenFetcher(httpClient, test.MockAuth, requestUrl)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())

	httpClient = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		return test.Responses.Empty()
	})

	token, err = client.ExportedAccessTokenFetcher(httpClient, test.MockAuth, nil)
	assert.Error(t, err)
	assert.Nil(t, token)
}
