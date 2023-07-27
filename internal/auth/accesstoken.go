package auth

import (
	"encoding/json"
	"time"

	"github.com/Xavier-Lam/go-wechat/caches"
)

const (
	BizAccessToken = "ak"

	DefaultAccessTokenExpiresIn = 7200
)

// AccessTokenManager is the credential manager to hold access token.
type AccessTokenManager = CredentialManager[AccessToken]

type accessTokenFetcher = func() (*AccessToken, error)

type AccessToken struct {
	accessToken string
	expiresIn   int
	createdAt   time.Time
}

// NewAccessToken creates a new `AccessToken` instance.
func NewAccessToken(accessToken string, expiresIn int) *AccessToken {
	if expiresIn <= 0 {
		expiresIn = DefaultAccessTokenExpiresIn
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

// accessTokenManager is an implement of the `auth.CredentialManager`
// which is used to manage access token credentials.
type accessTokenManager struct {
	auth  Auth
	cache caches.Cache
	fetch accessTokenFetcher
}

// NewAccessTokenManager creates a new instance of `auth.CredentialManager`
// to manage access token credentials.
func NewAccessTokenManager(auth Auth, cache caches.Cache, fetcher accessTokenFetcher) AccessTokenManager {
	return &accessTokenManager{
		auth:  auth,
		cache: cache,
		fetch: fetcher,
	}
}

func (cm *accessTokenManager) Get() (*AccessToken, error) {
	cachedValue, err := cm.get()
	if err == nil {
		return cachedValue, nil
	}

	return cm.Renew()
}

func (cm *accessTokenManager) Set(credential *AccessToken) error {
	return ErrNotSettable
}

func (cm *accessTokenManager) Renew() (*AccessToken, error) {
	if cm.fetch == nil {
		return nil, ErrNotRenewable
	}

	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.fetch()
	if err != nil {
		return nil, err
	}

	if cm.cache == nil {
		err = caches.ErrCacheNotSet
	} else {
		serializedToken, err := SerializeAccessToken(token)
		if err != nil {
			return nil, err
		}
		err = cm.cache.Set(
			cm.auth.GetAppId(),
			BizAccessToken,
			serializedToken,
			token.GetExpiresIn(),
		)
	}

	return token, err
}

func (cm *accessTokenManager) Delete() error {
	token, err := cm.get()
	if err != nil {
		return err
	}
	serializedToken, err := SerializeAccessToken(token)
	if err != nil {
		return err
	}
	return cm.cache.Delete(
		cm.auth.GetAppId(),
		BizAccessToken,
		serializedToken,
	)
}

func (cm *accessTokenManager) get() (*AccessToken, error) {
	if cm.cache == nil {
		return nil, caches.ErrCacheNotSet
	}

	cachedValue, err := cm.cache.Get(cm.auth.GetAppId(), BizAccessToken)
	if err != nil {
		return nil, err
	}

	token, err := DeserializeAccessToken(cachedValue)
	if err != nil {
		return nil, err
	}

	return token, nil
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
