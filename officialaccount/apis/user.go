package apis

import (
	"net/url"

	"github.com/Xavier-Lam/go-wechat/client"
)

type UserInfo struct {
	Subscribe      int    `json:"subscribe"`
	OpenId         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int    `json:"subscribe_time"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupId        int    `json:"groupid"`
	TagIdList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

type user struct {
	c client.WeChatClient
}

// User management
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
type User interface {
	// Obtaining Users' Basic Information
	// https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	GetInfo(openid string, lang string) (*UserInfo, error)
}

func newUser(client client.WeChatClient) User {
	return &user{c: client}
}

func (api *user) GetInfo(openid string, lang string) (*UserInfo, error) {
	q := url.Values{}
	q.Add("openid", openid)
	q.Add("lang", lang)
	// TODO: build query string
	endpoint := "/cgi-bin/user/info"
	resp, err := api.c.Get(endpoint+"?"+q.Encode(), true)
	if err != nil {
		return nil, err
	}
	userInfo := &UserInfo{}
	err = client.GetJson(resp, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
