package officialaccount

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/officialaccount/apis"
)

type Config struct {
	HttpClient        *http.Client             // Default Http client to send request
	Cache             caches.Cache             // Cache instance for managing tokens
	AccessTokenClient client.AccessTokenClient // The client used for request access token
	BaseApiUri        *url.URL                 // The endpoint to request an API, if full path is not given, default value is 'https://api.weixin.qq.com'
}

type App struct {
	Apis *apis.Apis

	Js js
}

func New(auth auth.Auth, conf Config) *App { // Set up base dependencies if not given
	c := client.New(auth, client.Config{
		HttpClient:        conf.HttpClient,
		Cache:             conf.Cache,
		AccessTokenClient: conf.AccessTokenClient,
		BaseApiUrl:        conf.BaseApiUri,
	})
	a := apis.NewApis(c)
	return &App{
		Apis: a,

		Js: *newJs(auth, a.Js, conf.Cache),
	}
}

func (a *App) GetAccessToken() (*client.Token, error) {
	return a.Apis.GetAccessToken()
}
