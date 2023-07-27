package examples

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
)

// To get access token from a customized URL
func CustomizeAccessTokenUrl() {
	var (
		appId     = os.Getenv("WECHAT_APP_ID")
		appSecret = os.Getenv("WECHAT_APP_SECRET")
	)

	auth := wechat.NewAuth(appId, appSecret)
	cache := caches.NewDummyCache()
	accessTokenUrl, _ := url.Parse("https://example.com/token")
	conf := wechat.OfficialAccountConfig{
		AccessTokenUrl: accessTokenUrl,
		Cache:          cache,
	}
	app := wechat.NewOfficeAccount(auth, conf)
	accessToken, err := app.GetAccessToken()
	if err == nil {
		fmt.Println(accessToken.GetAccessToken())
		fmt.Println(accessToken.GetExpiresIn())
	}
}

// To get access token by a customized way
func CustomizeAccessTokenFetcher() {
	var (
		appId     = os.Getenv("WECHAT_APP_ID")
		appSecret = os.Getenv("WECHAT_APP_SECRET")
	)

	auth := wechat.NewAuth(appId, appSecret)
	cache := caches.NewDummyCache()
	conf := wechat.OfficialAccountConfig{
		AccessTokenFetcher: func(client *http.Client, auth wechat.Auth, accessTokenUrl *url.URL) (*wechat.AccessToken, error) {
			return wechat.NewAccessToken("mock-access-token", 7200), nil
		},
		Cache: cache,
	}
	app := wechat.NewOfficeAccount(auth, conf)
	accessToken, err := app.GetAccessToken()
	if err == nil {
		fmt.Println(accessToken.GetAccessToken())
		fmt.Println(accessToken.GetExpiresIn())
	}
}
