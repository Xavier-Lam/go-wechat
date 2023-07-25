package auth_test

import (
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

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
	mockAuth := auth.New("", "")
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	atc := test.NewMockAccessTokenClient(oldToken)
	cm := auth.NewAccessTokenManager(atc, mockAuth, cache)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.(*auth.AccessToken).GetAccessToken())

	atc = test.NewMockAccessTokenClient(newToken)
	cm = auth.NewAccessTokenManager(atc, mockAuth, cache)

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

func TestAccessTokenManagerDelete(t *testing.T) {
	oldToken := "old"
	newToken := "token"
	mockAuth := auth.New("", "")

	cache := caches.NewDummyCache()
	atc := test.NewMockAccessTokenClient(oldToken)
	cm := auth.NewAccessTokenManager(atc, mockAuth, cache)

	err := cm.Delete()
	assert.Error(t, err)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, oldToken, token.(*auth.AccessToken).GetAccessToken())

	atc = test.NewMockAccessTokenClient(newToken)
	cm = auth.NewAccessTokenManager(atc, mockAuth, cache)

	err = cm.Delete()
	assert.NoError(t, err)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &auth.AccessToken{}, token)
	assert.Equal(t, newToken, token.(*auth.AccessToken).GetAccessToken())
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
