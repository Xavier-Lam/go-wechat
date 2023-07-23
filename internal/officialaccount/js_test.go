package officialaccount_test

import (
	"testing"
	"time"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount/apis"
	"github.com/stretchr/testify/assert"
)

type mockJsApi struct {
	ticket string
}

func newMockJsApi(ticket string) apis.Js {
	return &mockJsApi{
		ticket: ticket,
	}
}

func (api *mockJsApi) GetTicket() (*apis.JSTicket, error) {
	return &apis.JSTicket{
		Ticket:    api.ticket,
		ExpiresIn: 7200,
	}, nil
}

func TestJsGetTicket(t *testing.T) {
	oldTicket := "old"
	newTicket := "ticket"

	auth := wechat.NewAuth("app-id", "app-secret")
	cache := caches.NewDummyCache()
	js := officialaccount.NewJs(auth, newMockJsApi(oldTicket), cache)

	ticket, err := js.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, oldTicket, ticket)

	ticket, err = js.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, oldTicket, ticket)

	js = officialaccount.NewJs(auth, newMockJsApi(newTicket), cache)

	ticket, err = js.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, oldTicket, ticket)

	ticket, err = js.FetchTicket()
	assert.NoError(t, err)
	assert.Equal(t, newTicket, ticket)

	ticket, err = js.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, newTicket, ticket)
}

func TestJsConfig(t *testing.T) {
	nonceStr := "Wm3WZYTPz0wzccnW"
	ticket := "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg"
	timestamp := 1414587457
	url := "http://mp.weixin.qq.com?params=value"
	actualSignature := "0f9de62fce790f9a083d5c99e95740ceb90c27ed"
	auth := wechat.NewAuth("app-id", "app-secret")
	cache := caches.NewDummyCache()
	js := officialaccount.NewJs(auth, newMockJsApi(ticket), cache)

	jsConfig, err := js.GetJsConfig(url, officialaccount.JsConfig{
		Timestamp: timestamp,
		NonceStr:  nonceStr,
	})
	assert.NoError(t, err)
	assert.Equal(t, "app-id", jsConfig.AppId)
	assert.Equal(t, false, jsConfig.Debug)
	assert.Equal(t, []string{}, jsConfig.JsApiList)
	assert.Equal(t, nonceStr, jsConfig.NonceStr)
	assert.Equal(t, timestamp, jsConfig.Timestamp)
	assert.Equal(t, actualSignature, jsConfig.Signature)

	jsConfig, err = js.GetJsConfig(url, officialaccount.JsConfig{
		Debug:     true,
		JsApiList: []string{"a", "b"},
		Timestamp: timestamp,
		NonceStr:  nonceStr,
	})
	assert.NoError(t, err)
	assert.Equal(t, "app-id", jsConfig.AppId)
	assert.Equal(t, true, jsConfig.Debug)
	assert.Equal(t, []string{"a", "b"}, jsConfig.JsApiList)
	assert.Equal(t, nonceStr, jsConfig.NonceStr)
	assert.Equal(t, timestamp, jsConfig.Timestamp)
	assert.Equal(t, actualSignature, jsConfig.Signature)

	jsConfig, err = js.GetJsConfig(url, officialaccount.JsConfig{})
	assert.NoError(t, err)
	assert.Equal(t, "app-id", jsConfig.AppId)
	assert.Equal(t, false, jsConfig.Debug)
	assert.Equal(t, []string{}, jsConfig.JsApiList)
	assert.NotEmpty(t, jsConfig.NonceStr)
	assert.Equal(t, int(time.Now().Unix()), jsConfig.Timestamp)
	assert.NotEmpty(t, actualSignature)
}

func TestJsSign(t *testing.T) {
	nonceStr := "Wm3WZYTPz0wzccnW"
	ticket := "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg"
	timestamp := 1414587457
	url := "http://mp.weixin.qq.com?params=value"
	actualSignature := "0f9de62fce790f9a083d5c99e95740ceb90c27ed"

	auth := wechat.NewAuth("app-id", "app-secret")
	cache := caches.NewDummyCache()
	js := officialaccount.NewJs(auth, newMockJsApi(ticket), cache)

	signature, err := js.Sign(url, nonceStr, timestamp)
	assert.NoError(t, err)
	assert.Equal(t, actualSignature, signature)
}
