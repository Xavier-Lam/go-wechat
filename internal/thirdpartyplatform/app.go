package thirdpartyplatform

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform/apis"
)

type ThirdPartyPlatformAuth interface {
	auth.Auth
	GetTicket() (string, error)
}

type thirdPartyPlatformAuth struct {
	auth.Auth
	vtm VerifyTicketManager
}

func NewThirdPartyPlatformAuth(auth auth.Auth, vtm VerifyTicketManager) ThirdPartyPlatformAuth {
	return &thirdPartyPlatformAuth{auth, vtm}
}

func (a *thirdPartyPlatformAuth) GetTicket() (string, error) {
	ticket, err := a.vtm.Get()
	if err != nil {
		return "", err
	}
	return *ticket, nil
}

type Config struct {
	// BaseApiUrl is the base URL used for making API requests.
	// If not provided, the default value is 'https://api.weixin.qq.com'.
	BaseApiUrl *url.URL

	// Cache instance for managing tokens
	Cache caches.Cache

	// HttpClient is the default HTTP client used for sending requests.
	HttpClient *http.Client

	// ThirdPartyPlatformAccessTokenClient is used for request a latest access token when it is needed
	// This option should be left as the default value (nil), unless you want to customize the client
	// For example, if you want to request your access token from a different service than Tencent's.
	ThirdPartyPlatformAccessTokenClient auth.AccessTokenClient

	VerifyTicketManagerFactory VerifyTicketManagerFactory
}

type App struct {
	vtm VerifyTicketManager

	Apis *apis.Apis
}

func New(a auth.Auth, conf Config) *App {
	if conf.VerifyTicketManagerFactory == nil {
		conf.VerifyTicketManagerFactory = NewVerifyTicketManager
	}
	vtm := conf.VerifyTicketManagerFactory(a, conf.Cache, DefaultVerifyTicketExpiresIn)

	if conf.ThirdPartyPlatformAccessTokenClient == nil {
		conf.ThirdPartyPlatformAccessTokenClient = NewAccessTokenClient(conf.HttpClient, "")
	}

	tpa := NewThirdPartyPlatformAuth(a, vtm)

	c := client.New(tpa, client.Config{
		AccessTokenClient: conf.ThirdPartyPlatformAccessTokenClient,
		BaseApiUrl:        conf.BaseApiUrl,
		Cache:             conf.Cache,
		HttpClient:        conf.HttpClient,
	})

	api := apis.NewApis(c)
	return &App{
		vtm: vtm,

		Apis: api,
	}
}

func (a *App) GetTicket() (string, error) {
	ticket, err := a.vtm.Get()
	if err != nil {
		return "", err
	}
	return *ticket, nil
}

func (a *App) SetTicket(ticket string) error {
	return a.vtm.Set(&ticket)
}
