package miniprogram

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram/apis"
)

// Config can be extended in the future, it is not a reference to `client.Config`
// By using this reference, we can write less code for now.
// DO NOT use `client.Config` directly to avoid any potential future changes.
type Config = client.Config

type App struct {
	Apis *apis.Apis
}

func New(auth auth.Auth, conf Config) *App { // Set up base dependencies if not given
	c := client.New(auth, client.Config{
		CredentialManagerFactory: conf.CredentialManagerFactory,
		BaseApiUrl:               conf.BaseApiUrl,
		Cache:                    conf.Cache,
		HttpClient:               conf.HttpClient,
	})
	a := apis.NewApis(c)
	return &App{
		Apis: a,
	}
}

func (a *App) JsCode2Session(code string) (*apis.Session, error) {
	return a.Apis.Login.JsCode2Session(code)
}

func (a *App) GetAccessToken() (*auth.AccessToken, error) {
	return a.Apis.GetAccessToken()
}
