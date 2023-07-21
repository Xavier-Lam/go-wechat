package apis_test

import (
	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/officialaccount"
)

var (
	appID       = "mock-app-id"
	appSecret   = "mock-app-secret"
	accessToken = "mock-access-token"
	auth        = wechat.NewAuth(appID, appSecret)
)

func newMockOfficialAccount(handler test.RequestHandler) *officialaccount.OfficialAccount {
	return officialaccount.New(
		auth,
		officialaccount.Config{
			AccessTokenClient: test.NewMockAccessTokenClient(accessToken),
			HttpClient:        test.NewMockHttpClient(handler),
		},
	)
}
