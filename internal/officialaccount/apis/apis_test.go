package apis_test

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount"
	"github.com/Xavier-Lam/go-wechat/internal/test"
)

func newMockOfficialAccount(handler test.RequestHandler) *officialaccount.App {
	return officialaccount.New(
		test.MockAuth,
		officialaccount.Config{
			AccessTokenFetcher: func(c *http.Client, a auth.Auth, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
				return auth.NewAccessToken(test.AccessToken, 0), nil
			},
			HttpClient: test.NewMockHttpClient(handler),
		},
	)
}
