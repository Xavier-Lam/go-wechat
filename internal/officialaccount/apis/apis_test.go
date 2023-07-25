package apis_test

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount"
	"github.com/Xavier-Lam/go-wechat/internal/test"
)

var (
	appID       = "mock-app-id"
	appSecret   = "mock-app-secret"
	accessToken = "mock-access-token"
	mockAuth    = auth.New(appID, appSecret)
)

func newMockOfficialAccount(handler test.RequestHandler) *officialaccount.App {
	return officialaccount.New(
		mockAuth,
		officialaccount.Config{
			AccessTokenClient: test.NewMockAccessTokenClient(accessToken),
			HttpClient:        test.NewMockHttpClient(handler),
		},
	)
}
