package auth_test

import (
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestTokenGetExpires(t *testing.T) {
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

func TestTokenSerialize(t *testing.T) {
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
