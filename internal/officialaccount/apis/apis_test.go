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
	mockAuth    = auth.NewAuth(appID, appSecret)
)

func newMockOfficialAccount(handler test.RequestHandler) *officialaccount.App {
	return officialaccount.New(
		mockAuth,
		officialaccount.Config{
			CredentialManagerFactory: test.MockAccessTokenCredentialManagerFactoryProvider(accessToken),
			HttpClient:               test.NewMockHttpClient(handler),
		},
	)
}
