package apis

import (
	"net/url"

	"github.com/Xavier-Lam/go-wechat/client"
)

type Session struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type login struct {
	c client.WeChatClient
}

// User management
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
type Login interface {
	// Obtaining Users' Basic Information
	// https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
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
	return rv, nil
}
