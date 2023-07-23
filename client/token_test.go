package client_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTokenGetAccessToken(t *testing.T) {
	auth := wechat.NewAuth("app-id", "app-secret")

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
	c := client.NewAccessTokenClient(url, client.NewAuthCredentialManager(auth), httpClient)

	token, err := c.GetAccessToken()

	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*7200), token.GetExpiresAt(), time.Millisecond*50)
}

func TestTokenGetExpires(t *testing.T) {
	token := client.NewToken("access_token", 2)

	assert.Equal(t, 2, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*2), token.GetExpiresAt(), time.Millisecond*50)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*1), token.GetExpiresAt(), time.Millisecond*50)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*0), token.GetExpiresAt(), time.Millisecond*50)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, token.GetExpiresIn())
	assert.WithinDuration(t, time.Now().Add(time.Second*-1), token.GetExpiresAt(), time.Millisecond*50)
}

func TestTokenSerialize(t *testing.T) {
	token := client.NewToken("access_token", 2)
	// Serialize the token
	bytes, err := client.SerializeToken(token)
	assert.NoError(t, err)

	// Deserialize the bytes back into a token
	deserializedToken, err := client.DeserializeToken(bytes)
	assert.NoError(t, err)

	assert.Equal(t, token.GetAccessToken(), deserializedToken.GetAccessToken())
	assert.WithinDuration(t, token.GetExpiresAt(), deserializedToken.GetExpiresAt(), time.Millisecond*50)
	assert.Equal(t, token.GetExpiresIn(), deserializedToken.GetExpiresIn())
}
