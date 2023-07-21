package apis

import "github.com/Xavier-Lam/go-wechat/client"

type Apis struct {
	client.WeChatClient

	Js   Js
	User User
}

func NewApis(c client.WeChatClient) *Apis {
	return &Apis{
		c,

		newJs(c),
		newUser(c),
	}
}
