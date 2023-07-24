package auth

import (
	"encoding/json"
	"time"
)

const DefaultTokenExpiresIn = 7200

type AccessToken struct {
	accessToken string
	expiresIn   int
	createdAt   time.Time
}

// NewAccessToken creates a new `AccessToken` instance.
func NewAccessToken(accessToken string, expiresIn int) *AccessToken {
	if expiresIn <= 0 {
		expiresIn = DefaultTokenExpiresIn
	}
	return &AccessToken{
		accessToken: accessToken,
		expiresIn:   expiresIn,
		createdAt:   time.Now(),
	}
}

// GetAccessToken returns the access token value.
func (t *AccessToken) GetAccessToken() string {
	return t.accessToken
}

// GetExpiresIn returns the remaining time until the access token expires in seconds.
func (t *AccessToken) GetExpiresIn() int {
	timeDiff := time.Since(t.createdAt)
	timeEscaped := int(timeDiff.Seconds())
	if timeEscaped >= t.expiresIn {
		return 0
	}
	return t.expiresIn - timeEscaped
}

// GetExpiresAt returns the time when the access token will expire.
func (t *AccessToken) GetExpiresAt() time.Time {
	timeDiff := time.Duration(t.expiresIn) * time.Second
	return t.createdAt.Add(timeDiff)
}

type accessToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	CreatedAt   time.Time `json:"created_at"`
}

func SerializeAccessToken(token *AccessToken) ([]byte, error) {
	timeDiff := -time.Duration(time.Second * time.Duration(token.GetExpiresIn()))
	data := &accessToken{
		AccessToken: token.GetAccessToken(),
		ExpiresIn:   token.GetExpiresIn(),
		CreatedAt:   token.GetExpiresAt().Add(timeDiff),
	}
	return json.Marshal(data)
}

func DeserializeAccessToken(bytes []byte) (*AccessToken, error) {
	data := &accessToken{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		accessToken: data.AccessToken,
		expiresIn:   data.ExpiresIn,
		createdAt:   data.CreatedAt,
	}, nil
}
