package apis_test

import (
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestLoginJsCode2Session(t *testing.T) {
	data := `{
		"openid": "xxxxxx",
		"session_key": "xxxxx",
		"unionid": "xxxxx",
		"errcode": 0,
		"errmsg": "xxxxx"
	}`

	app := newMockMiniProgram(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, "https://api.weixin.qq.com/sns/jscode2session", req.URL)
		assert.Equal(t, test.AppId, req.URL.Query().Get("appid"))
		assert.Equal(t, test.AppSecret, req.URL.Query().Get("secret"))
		assert.Equal(t, "TEST_CODE", req.URL.Query().Get("js_code"))
		assert.Equal(t, "authorization_code", req.URL.Query().Get("grant_type"))

		return test.Responses.Json(data)
	})

	session, err := app.Apis.Login.JsCode2Session("TEST_CODE")
	assert.NoError(t, err)
	assert.Equal(t, "xxxxxx", session.OpenId)
	assert.Equal(t, "xxxxx", session.SessionKey)
	assert.Equal(t, "xxxxx", session.UnionId)
}
