package client_test

import (
	"testing"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestWeChatClientCredential(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	akc := test.NewMockAccessTokenClient(oldToken)
	cm := client.NewWeChatAccessTokenCredentialManager(auth, cache, akc)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, oldToken, token.(*client.Token).GetAccessToken())

	akc = test.NewMockAccessTokenClient(newToken)
	cm = client.NewWeChatAccessTokenCredentialManager(auth, cache, akc)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, oldToken, token.(*client.Token).GetAccessToken())

	token, err = cm.Renew()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, newToken, token.(*client.Token).GetAccessToken())

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, newToken, token.(*client.Token).GetAccessToken())
}

func TestWeChatClientCredentialDelete(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	akc := test.NewMockAccessTokenClient(oldToken)
	cm := client.NewWeChatAccessTokenCredentialManager(auth, cache, akc)

	err := cm.Delete()
	assert.Error(t, err)

	token, err := cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, oldToken, token.(*client.Token).GetAccessToken())

	akc = test.NewMockAccessTokenClient(newToken)
	cm = client.NewWeChatAccessTokenCredentialManager(auth, cache, akc)

	err = cm.Delete()
	assert.NoError(t, err)

	token, err = cm.Get()
	assert.NoError(t, err)
	assert.IsType(t, &client.Token{}, token)
	assert.Equal(t, newToken, token.(*client.Token).GetAccessToken())
}
