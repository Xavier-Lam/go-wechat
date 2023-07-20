package apis_test

import (
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestUserGetInfo(t *testing.T) {
	openid := "o6_bmjrPTlm6_2sgVt7hMZOPfL2M"
	lang := "zh_CN"
	data := `{
		"subscribe": 1, 
		"openid": "o6_bmjrPTlm6_2sgVt7hMZOPfL2M", 
		"language": "zh_CN", 
		"subscribe_time": 1382694957,
		"unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
		"remark": "",
		"groupid": 0,
		"tagid_list":[128,2],
		"subscribe_scene": "ADD_SCENE_QR_CODE",
		"qr_scene": 98765,
		"qr_scene_str": ""
	}`

	app := newMockOfficialAccount(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		test.AssertEndpointEqual(t, "https://api.weixin.qq.com/cgi-bin/user/info", req.URL)
		assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))
		assert.Equal(t, openid, req.URL.Query().Get("openid"))
		assert.Equal(t, lang, req.URL.Query().Get("lang"))

		return test.Responses.Json(data)
	})

	userInfo, err := app.Apis.User.GetInfo(openid, lang)
	assert.NoError(t, err)
	assert.Equal(t, 1, userInfo.Subscribe)
	assert.Equal(t, openid, userInfo.OpenId)
	assert.Equal(t, lang, userInfo.Language)
	assert.Equal(t, 1382694957, userInfo.SubscribeTime)
	assert.Equal(t, "o6_bmasdasdsad6_2sgVt7hMZOPfL", userInfo.UnionId)
	assert.Equal(t, "", userInfo.Remark)
	assert.Equal(t, 0, userInfo.GroupId)
	assert.Equal(t, []int{128, 2}, userInfo.TagIdList)
	assert.Equal(t, "ADD_SCENE_QR_CODE", userInfo.SubscribeScene)
	assert.Equal(t, 98765, userInfo.QrScene)
	assert.Equal(t, "", userInfo.QrSceneStr)
}
