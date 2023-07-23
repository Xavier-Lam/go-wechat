package client_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTokenGetAccessToken(t *testing.T) {
	a := wechat.NewAuth("app-id", "app-secret")

	httpClient := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, client.DefaultAccessTokenUri, req.URL)
		assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
		assert.Equal(t, "app-id", req.URL.Query().Get("appid"))
		assert.Equal(t, "app-secret", req.URL.Query().Get("secret"))

		return test.Responses.Json(`{"access_token": "access-token", "expires_in": 7200}`)
	})
	url, _ := url.Parse(client.DefaultAccessTokenUri)
	c := client.NewAccessTokenClient(url, a, httpClient)

	token, err := c.GetAccessToken()

	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*7200), token.GetExpiresAt(), time.Millisecond*50)
}

func TestWeChatAccessTokenCredential(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	akc := test.NewMockAccessTokenClient(oldToken)
	cm := client.NewWeChatAccessTokenCredentialManager(mockAuth, cache, akc)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.(*auth.AccessToken).GetAccessToken())

	akc = test.NewMockAccessTokenClient(newToken)
	cm = client.NewWeChatAccessTokenCredentialManager(mockAuth, cache, akc)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.(*auth.AccessToken).GetAccessToken())

	token, err = cm.Renew()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.(*auth.AccessToken).GetAccessToken())

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.(*auth.AccessToken).GetAccessToken())
}

func TestWeChatAccessTokenCredentialDelete(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	akc := test.NewMockAccessTokenClient(oldToken)
	cm := client.NewWeChatAccessTokenCredentialManager(mockAuth, cache, akc)

	err := cm.Delete()
	assert.Error(t, err)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.(*auth.AccessToken).GetAccessToken())

	akc = test.NewMockAccessTokenClient(newToken)
	cm = client.NewWeChatAccessTokenCredentialManager(mockAuth, cache, akc)

	err = cm.Delete()
	assert.NoError(t, err)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.(*auth.AccessToken).GetAccessToken())
}
