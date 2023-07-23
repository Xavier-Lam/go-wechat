package apis

import "github.com/Xavier-Lam/go-wechat/internal/client"

type Apis struct {
	client.WeChatClient

	Login Login
}

func NewApis(c client.WeChatClient) *Apis {
	return &Apis{
		c,

		newLogin(c),
	}
}
