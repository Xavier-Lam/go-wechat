package apis_test

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
)

const mockTicket = "ticket"

type mockVerifyTicketManager struct {
	ticket string
}

func (cm *mockVerifyTicketManager) Get() (*string, error) {
	return &cm.ticket, nil
}

func (cm *mockVerifyTicketManager) Set(credential *string) error {
	return auth.ErrNotSettable
}

func (cm *mockVerifyTicketManager) Renew() (*string, error) {
	return nil, auth.ErrNotRenewable
}

func (cm *mockVerifyTicketManager) Delete() error {
	return auth.ErrNotDeletable
}

func newMockThirdPartyPlatform(handler test.RequestHandler) *thirdpartyplatform.App {
	return thirdpartyplatform.New(
		test.MockAuth,
		thirdpartyplatform.Config{
			AccessTokenFetcher: func(c *http.Client, a auth.Auth, t string, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
				return auth.NewAccessToken(test.AccessToken, 0), nil
			},
			Cache:      caches.NewDummyCache(),
			HttpClient: test.NewMockHttpClient(handler),
		},
	)
}
