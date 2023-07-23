package apis

import (
	"github.com/Xavier-Lam/go-wechat/internal/client"
)

type JSTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// JS SDK
// https://developers.weixin.qq.com/doc/offiaccount/en/OA_Web_Apps/JS-SDK.html
type Js interface {
	// Get the latest validate ticket
	// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
	GetTicket() (*JSTicket, error)
}

type js struct {
	c client.WeChatClient
}

func newJs(c client.WeChatClient) Js {
	return &js{c: c}
}

func (api *js) GetTicket() (*JSTicket, error) {
	resp, err := api.c.Get("/cgi-bin/ticket/getticket?type=wx_card", true)
	if err != nil {
		return nil, err
	}
	ticket := &JSTicket{}
	err = client.GetJson(resp, ticket)
	if err != nil {
		return nil, err
	}
	if ticket.ExpiresIn <= 0 {
		ticket.ExpiresIn = 7200
	}
	return ticket, nil
}
