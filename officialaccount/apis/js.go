package apis

import (
	"fmt"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
)

type JSTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// JS SDK
// https://developers.weixin.qq.com/doc/offiaccount/en/OA_Web_Apps/JS-SDK.html
type Js interface {
	// Get the latest validate ticket (obtaining from cache first)
	// It may return an error along with the ticket if there is no `Cache` set up.
	// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
	GetTicket() (*JSTicket, error)
	// Obtaining api_ticket from server side
	// It may return an error along with the ticket if there is no `Cache` set up.
	// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
	FetchTicket() (*JSTicket, error)
}

type js struct {
	c     client.WeChatClient
	cache caches.Cache
}

func NewJs(c client.WeChatClient, cache caches.Cache) Js {
	return &js{
		c:     c,
		cache: cache,
	}
}

func (api *js) GetTicket() (*JSTicket, error) {
	if api.cache != nil {
		cachedValue, err := api.cache.Get(api.c.GetAuth().GetAppId(), caches.CacheBizJSTicket)
		if err == nil {
			if ticket, ok := cachedValue.(*JSTicket); ok {
				return ticket, nil
			}
		}
	}

	return api.FetchTicket()
}

func (api *js) FetchTicket() (*JSTicket, error) {
	ticket, err := api.getTicket()
	if err != nil {
		return nil, err
	}

	if api.cache == nil {
		err = fmt.Errorf("Cache is not set")
	} else {
		err = api.cache.Set(api.c.GetAuth().GetAppId(), caches.CacheBizJSTicket, ticket, ticket.ExpiresIn)
	}

	return ticket, err
}

func (api *js) getTicket() (*JSTicket, error) {
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
