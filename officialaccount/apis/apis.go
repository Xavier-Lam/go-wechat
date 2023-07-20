package apis

import "github.com/Xavier-Lam/go-wechat/client"

type Apis struct {
	c client.WeChatClient

	Js   Js
	User User
}

func NewApis(c client.WeChatClient) *Apis {
	return &Apis{
		c: c,

		// Js: NewJs(c),
		User: newUser(c),
	}
}
