package auth_test

import (
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

var mockAccessTokenFetcher = func(token string) func() (*auth.AccessToken, error) {
	return func() (*auth.AccessToken, error) {
		return auth.NewAccessToken(token, 0), nil
	}
}

func TestAccessTokenGetExpires(t *testing.T) {
	token := auth.NewAccessToken("access_token", 2)

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

func TestAccessTokenManager(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	cm := auth.NewAccessTokenManager(test.MockAuth, cache, mockAccessTokenFetcher(oldToken))

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.GetAccessToken())

	cm = auth.NewAccessTokenManager(test.MockAuth, cache, mockAccessTokenFetcher(newToken))

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.GetAccessToken())

	token, err = cm.Renew()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.GetAccessToken())

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.GetAccessToken())
}

func TestAccessTokenManagerDelete(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	cm := auth.NewAccessTokenManager(test.MockAuth, cache, mockAccessTokenFetcher(oldToken))

	err := cm.Delete()
	assert.Error(t, err)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.GetAccessToken())

	cm = auth.NewAccessTokenManager(test.MockAuth, cache, mockAccessTokenFetcher(newToken))

	err = cm.Delete()
	assert.NoError(t, err)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.GetAccessToken())
}

func TestAccessTokenSerialize(t *testing.T) {
	token := auth.NewAccessToken("access_token", 2)
	// Serialize the token
	bytes, err := auth.SerializeAccessToken(token)
	assert.NoError(t, err)

	// Deserialize the bytes back into a token
	deserializedToken, err := auth.DeserializeAccessToken(bytes)
	assert.NoError(t, err)

	assert.Equal(t, token.GetAccessToken(), deserializedToken.GetAccessToken())
	assert.WithinDuration(t, token.GetExpiresAt(), deserializedToken.GetExpiresAt(), time.Millisecond*50)
	assert.Equal(t, token.GetExpiresIn(), deserializedToken.GetExpiresIn())
}
