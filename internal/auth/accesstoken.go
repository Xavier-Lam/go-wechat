package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Xavier-Lam/go-wechat/caches"
)

const DefaultAccessTokenExpiresIn = 7200

// AccessTokenClient is an client to request the newest access token
type AccessTokenClient interface {
	PrepareRequest(auth Auth) (*http.Request, error)
	SendRequest(auth Auth, req *http.Request) (*http.Response, error)
	HandleResponse(auth Auth, resp *http.Response, req *http.Request) (*AccessToken, error)
}

// AccessTokenManager is the credential manager to hold access token.
type AccessTokenManager = CredentialManager[AccessToken]

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
	atc   AccessTokenClient
	auth  Auth
	cache caches.Cache
}

// NewAccessTokenManager creates a new instance of `auth.CredentialManager`
// to manage access token credentials.
func NewAccessTokenManager(atc AccessTokenClient, auth Auth, cache caches.Cache) AccessTokenManager {
	return &accessTokenManager{
		atc:   atc,
		auth:  auth,
		cache: cache,
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
	return errors.New("not settable")
}

func (cm *accessTokenManager) Renew() (*AccessToken, error) {
	cm.Delete()

	// TODO: prevent concurrent fetching
	token, err := cm.getAccessToken()
	if err != nil {
		return nil, err
	}

	if cm.cache == nil {
		err = fmt.Errorf("cache is not set")
	} else {
		serializedToken, err := SerializeAccessToken(token)
		if err != nil {
			return nil, err
		}
		err = cm.cache.Set(
			cm.auth.GetAppId(),
			caches.BizAccessToken,
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
		caches.BizAccessToken,
		serializedToken,
	)
}

func (cm *accessTokenManager) get() (*AccessToken, error) {
	if cm.cache == nil {
		return nil, ErrCacheNotSet
	}

	cachedValue, err := cm.cache.Get(cm.auth.GetAppId(), caches.BizAccessToken)
	if err != nil {
		return nil, err
	}

	token, err := DeserializeAccessToken(cachedValue)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (cm *accessTokenManager) getAccessToken() (*AccessToken, error) {
	req, err := cm.atc.PrepareRequest(cm.auth)
	if err != nil {
		return nil, err
	}

	resp, err := cm.atc.SendRequest(cm.auth, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return cm.atc.HandleResponse(cm.auth, resp, req)
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
