package apis_test

import (
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
)

func newMockThirdPartyPlatform(handler test.RequestHandler) *thirdpartyplatform.App {
	return thirdpartyplatform.New(
		test.MockAuth,
		thirdpartyplatform.Config{
			ThirdPartyPlatformAccessTokenClient: test.NewMockAccessTokenClient(test.AccessToken),
			HttpClient:                          test.NewMockHttpClient(handler),
		},
	)
}
