package apis_test

import (
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestJsGetTicket(t *testing.T) {
	data := `{
		"errcode":0,
		"errmsg":"ok",
		"ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKdvsdshFKA",
		"expires_in":7200
	}`

	app := newMockOfficialAccount(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, "https://api.weixin.qq.com/cgi-bin/ticket/getticket", req.URL)
		assert.Equal(t, test.AccessToken, req.URL.Query().Get("access_token"))
		assert.Equal(t, "wx_card", req.URL.Query().Get("type"))

		return test.Responses.Json(data)
	})

	resp, err := app.Apis.Js.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, "bxLdikRXVbTPdHSM05e5u5sUoXNKdvsdshFKA", resp.Ticket)
	assert.Equal(t, 7200, resp.ExpiresIn)
}
