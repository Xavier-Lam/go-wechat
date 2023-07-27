package thirdpartyplatform

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform/apis"
)

type Config struct {
	// AccessTokenFetcher is a callback function to return the latest access token
	// The default implement should be suitable for most case, override only when
	// you want to customize the way you make request.
	// For example, if you want to request to a service rather than Tencent's.
	AccessTokenFetcher AccessTokenFetcher

	// AccessTokenUrl is the url AccessTokenFetcher tries to fetch the latest access token.
	// This URL will be passed to the AccessTokenFetcher callback.
	// If not provided, the default value is 'https://api.weixin.qq.com/cgi-bin/component/api_component_token'.
	AccessTokenUrl *url.URL

	// BaseApiUrl is the base URL used for making API requests.
	// If not provided, the default value is 'https://api.weixin.qq.com'.
	BaseApiUrl *url.URL

	// Cache instance for managing tokens
	Cache caches.Cache

	// HttpClient is the default HTTP client used for sending requests.
	HttpClient *http.Client

	VerifyTicketManagerFactory VerifyTicketManagerFactory
}

type App struct {
	vtm VerifyTicketManager

	Apis *apis.Apis
}

func New(a auth.Auth, conf Config) *App {
	if conf.AccessTokenFetcher == nil {
		conf.AccessTokenFetcher = accessTokenFetcher
	}

	if conf.AccessTokenUrl == nil {
		conf.AccessTokenUrl, _ = url.Parse(DefaultAccessTokenUrl)
	}

	if conf.VerifyTicketManagerFactory == nil {
		conf.VerifyTicketManagerFactory = NewVerifyTicketManager
	}
	vtm := conf.VerifyTicketManagerFactory(a, conf.Cache, DefaultVerifyTicketExpiresIn)

	fetcher := func(c *http.Client, a auth.Auth, accessTokenUrl *url.URL) (*auth.AccessToken, error) {
		var ticketStr string
		ticket, err := vtm.Get()
		if err == nil {
			// Let fetcher to determine how to handle an empty ticket
			ticketStr = *ticket
		} else {
			ticketStr = ""
		}
		return conf.AccessTokenFetcher(c, a, ticketStr, accessTokenUrl)
	}

	c := client.New(a, client.Config{
		AccessTokenFetcher: fetcher,
		AccessTokenUrl:     conf.AccessTokenUrl,
		BaseApiUrl:         conf.BaseApiUrl,
		Cache:              conf.Cache,
		HttpClient:         conf.HttpClient,
	})

	api := apis.NewApis(c)
	return &App{
		vtm: vtm,

		Apis: api,
	}
}

func (a *App) GetAccessToken() (*auth.AccessToken, error) {
	return a.Apis.GetAccessToken()
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
