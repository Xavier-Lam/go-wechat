package client_test

import (
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestPrepareRequest(t *testing.T) {
	c := client.NewAccessTokenClient(nil, "")

	req, err := c.PrepareRequest(mockAuth)
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	test.AssertEndpointEqual(t, client.DefaultAccessTokenUrl, req.URL)
	assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
	assert.Equal(t, appID, req.URL.Query().Get("appid"))
	assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

	requestUrl := "https://example.com/test"
	c = client.NewAccessTokenClient(nil, requestUrl)

	req, err = c.PrepareRequest(mockAuth)
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	test.AssertEndpointEqual(t, requestUrl, req.URL)
	assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
	assert.Equal(t, appID, req.URL.Query().Get("appid"))
	assert.Equal(t, appSecret, req.URL.Query().Get("secret"))
}

func TestSendRequest(t *testing.T) {
	httpClient := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, client.DefaultAccessTokenUrl, req.URL)
		assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
		assert.Equal(t, appID, req.URL.Query().Get("appid"))
		assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

		return test.Responses.Empty()
	})
	c := client.NewAccessTokenClient(httpClient, "")

	req, err := c.PrepareRequest(mockAuth)
	assert.NoError(t, err)
	resp, err := c.SendRequest(mockAuth, req)

	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestHandleResponse(t *testing.T) {
	resp, _ := test.Responses.Json(`{"access_token": "access-token", "expires_in": 7200}`)
	client := client.NewAccessTokenClient(nil, "")

	token, err := client.HandleResponse(mockAuth, resp, &http.Request{})

	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())
}
