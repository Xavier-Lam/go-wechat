package apis

import (
	"net/url"

	"github.com/Xavier-Lam/go-wechat/internal/client"
)

type Session struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

type login struct {
	c client.WeChatClient
}

// Login
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
type Login interface {
	// Get miniprogram session
	// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	JsCode2Session(code string) (*Session, error)
}

func newLogin(c client.WeChatClient) Login {
	return &login{c: c}
}

func (api *login) JsCode2Session(code string) (*Session, error) {
	q := url.Values{}
	q.Add("appid", api.c.GetAuth().GetAppId())
	q.Add("secret", api.c.GetAuth().GetAppSecret())
	q.Add("js_code", code)
	q.Add("grant_type", "authorization_code")
	endpoint := "/sns/jscode2session"
	resp, err := api.c.Get(endpoint+"?"+q.Encode(), true)
	if err != nil {
		return nil, err
	}
	rv := &Session{}
	err = client.GetJson(resp, rv)
	if err != nil {
		return nil, err
	}

	if rv.OpenId == "" || rv.SessionKey == "" {
		return nil, client.ErrInvalidResponse
	}

	return rv, nil
}
