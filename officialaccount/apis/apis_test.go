package apis_test

import (
	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/officialaccount"
)

var (
	appID       = "mock-app-id"
	appSecret   = "mock-app-secret"
	accessToken = "mock-access-token"
	auth        = wechat.NewAuth(appID, appSecret)
)

type dummyCache struct {
	token string
}

func newDummyCache(token string) caches.Cache {
	return &dummyCache{token: token}
}

func (c *dummyCache) Get(appId string, biz string) (interface{}, error) {
	return client.NewToken(c.token, 0), nil
}

func (c *dummyCache) Set(appId string, biz string, value interface{}, expiresIn int) error {
	return nil
}

func newMockOfficialAccount(handler test.RequestHandler) *officialaccount.OfficialAccount {
	return officialaccount.New(
		auth,
		&officialaccount.Config{
			Cache:      newDummyCache(accessToken),
			HttpClient: test.NewMockHttpClient(handler),
		},
	)
}
