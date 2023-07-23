package apis_test

import (
	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram"
	"github.com/Xavier-Lam/go-wechat/internal/test"
)

var (
	appID       = "mock-app-id"
	appSecret   = "mock-app-secret"
	accessToken = "mock-access-token"
	auth        = wechat.NewAuth(appID, appSecret)
)

func newMockMiniProgram(handler test.RequestHandler) *miniprogram.App {
	return miniprogram.New(
		auth,
		miniprogram.Config{
			AccessTokenClient: test.NewMockAccessTokenClient(accessToken),
			HttpClient:        test.NewMockHttpClient(handler),
		},
	)
}
