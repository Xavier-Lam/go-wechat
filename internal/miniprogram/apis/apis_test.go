package apis_test

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram"
	"github.com/Xavier-Lam/go-wechat/internal/test"
)

func newMockMiniProgram(handler test.RequestHandler) *miniprogram.App {
	return miniprogram.New(
		test.MockAuth,
		miniprogram.Config{
			AccessTokenFetcher: func(client *http.Client, a auth.Auth, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
				return auth.NewAccessToken(test.AccessToken, 0), nil
			},
			HttpClient: test.NewMockHttpClient(handler),
		},
	)
}
