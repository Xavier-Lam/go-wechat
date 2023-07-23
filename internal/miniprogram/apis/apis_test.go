package apis_test

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram"
	"github.com/Xavier-Lam/go-wechat/internal/test"
)

var (
	appID       = "mock-app-id"
	appSecret   = "mock-app-secret"
	accessToken = "mock-access-token"
	mockAuth    = auth.NewAuth(appID, appSecret)
)

func newMockMiniProgram(handler test.RequestHandler) *miniprogram.App {
	return miniprogram.New(
		mockAuth,
		miniprogram.Config{
			CredentialManagerFactory: test.MockAccessTokenCredentialManagerFactoryProvider(accessToken),
			HttpClient:               test.NewMockHttpClient(handler),
		},
	)
}
