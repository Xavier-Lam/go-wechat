package apis

import (
	"github.com/Xavier-Lam/go-wechat/internal/client"
)

type PreAuthCode struct {
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

type ticketToken struct {
	c client.WeChatClient
}

// Ticket Token Management
type TicketToken interface {
	// Create Pre-Authorization Code
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
	GetPreAuthCode() (*PreAuthCode, error)
}

func newTicketToken(c client.WeChatClient) TicketToken {
	return &ticketToken{c: c}
}

func (api *ticketToken) GetPreAuthCode() (*PreAuthCode, error) {
	data := map[string]interface{}{
		"component_appid": api.c.GetAuth().GetAppId(),
	}

	resp, err := api.c.PostJson("/cgi-bin/component/api_create_preauthcode", data, true)
	if err != nil {
		return nil, err
	}

	preAuthCode := &PreAuthCode{}
	err = client.GetJson(resp, preAuthCode)
	if err != nil {
		return nil, err
	}

	if preAuthCode.PreAuthCode == "" {
		return nil, client.ErrInvalidResponse
	}

	return preAuthCode, nil
}
