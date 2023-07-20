package officialaccount

import (
	"net/http"
	"net/url"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/officialaccount/apis"
)

type Config struct {
	HttpClient     client.HttpClient // Default Http client to send request
	Cache          caches.Cache      // Cache instance for managing tokens
	AccessTokenUri *url.URL          // The endpoint to request a new token, default value is 'https://api.weixin.qq.com/cgi-bin/token'
	BaseApiUri     *url.URL          // The endpoint to request an API, if full path is not given, default value is 'https://api.weixin.qq.com'
}

type OfficialAccount struct {
	Apis *apis.Apis
}

func New(auth wechat.Auth, conf *Config) *OfficialAccount { // Set up base dependencies if not given
	if conf == nil {
		conf = &Config{}
	}
	if conf.AccessTokenUri == nil {
		conf.AccessTokenUri, _ = url.Parse(client.DefaultAccessTokenUri)
	}
	if conf.HttpClient == nil {
		conf.HttpClient = &http.Client{}
	}

	c := client.New(auth, &client.Config{
		HttpClient:     conf.HttpClient,
		Cache:          conf.Cache,
		AccessTokenUri: conf.AccessTokenUri,
		BaseApiUri:     conf.BaseApiUri,
	})
	return &OfficialAccount{
		Apis: apis.NewApis(c),
	}
}
