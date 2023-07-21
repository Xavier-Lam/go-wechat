package officialaccount

import (
	"net/url"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/officialaccount/apis"
)

type Config struct {
	HttpClient        client.HttpClient        // Default Http client to send request
	Cache             caches.Cache             // Cache instance for managing tokens
	AccessTokenClient client.AccessTokenClient // The client used for request access token
	BaseApiUri        *url.URL                 // The endpoint to request an API, if full path is not given, default value is 'https://api.weixin.qq.com'
}

type OfficialAccount struct {
	Apis *apis.Apis

	Js js
}

func New(auth wechat.Auth, conf Config) *OfficialAccount { // Set up base dependencies if not given
	c := client.New(auth, client.Config{
		HttpClient:        conf.HttpClient,
		Cache:             conf.Cache,
		AccessTokenClient: conf.AccessTokenClient,
		BaseApiUri:        conf.BaseApiUri,
	})
	a := apis.NewApis(c)
	return &OfficialAccount{
		Apis: a,

		Js: *newJs(auth, a.Js, conf.Cache),
	}
}
