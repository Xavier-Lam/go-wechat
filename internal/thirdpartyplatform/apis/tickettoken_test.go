package apis_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTicketTokenGetPreAuthCode(t *testing.T) {
	data := `{
		"pre_auth_code": "PRE_AUTH_CODE_VALUE",
		"expires_in": 7200
	}`

	app := newMockThirdPartyPlatform(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "POST", req.Method)
		test.AssertEndpointEqual(t, "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode", req.URL)
		assert.Equal(t, test.AccessToken, req.URL.Query().Get("access_token"))

		var body struct {
			AppId string `json:"component_appid"`
		}
		err := json.NewDecoder(req.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, test.AppId, body.AppId)

		return test.Responses.Json(data)
	})

	preAuthCode, err := app.Apis.TicketToken.GetPreAuthCode()
	assert.NoError(t, err)
	assert.Equal(t, "PRE_AUTH_CODE_VALUE", preAuthCode.PreAuthCode)
	assert.Equal(t, 7200, preAuthCode.ExpiresIn)
}
