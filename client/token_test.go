package client_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/stretchr/testify/assert"
)

func TestTokenGetAccessToken(t *testing.T) {
	auth := wechat.NewAuth("app-id", "app-secret")

	httpClient := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		assertEndpointEqual(t, client.DefaultAccessTokenUri, req.URL)
		assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
		assert.Equal(t, "app-id", req.URL.Query().Get("appid"))
		assert.Equal(t, "app-secret", req.URL.Query().Get("secret"))

		return createJsonResponse(`{"access_token": "access-token", "expires_in": 7200}`)
	})
	url, _ := url.Parse(client.DefaultAccessTokenUri)
	client := client.NewAccessTokenClient(url, httpClient)

	token, err := client.GetAccessToken(auth)

	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*7200), token.GetExpiresAt(), time.Second)
}

func TestTokenGetExpires(t *testing.T) {
	token := client.NewToken("access_token", 2)

	assert.Equal(t, 2, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*2), token.GetExpiresAt(), time.Second)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*1), token.GetExpiresAt(), time.Second)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*0), token.GetExpiresAt(), time.Second)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*-1), token.GetExpiresAt(), time.Second)
}
