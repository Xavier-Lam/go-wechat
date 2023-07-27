package examples

import (
	"fmt"
	"os"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
)

// Get miniprogram session
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
func GetSession() {
	var (
		appId     = os.Getenv("WECHAT_APP_ID")
		appSecret = os.Getenv("WECHAT_APP_SECRET")
	)

	code := "code"

	auth := wechat.NewAuth(appId, appSecret)
	cache := caches.NewDummyCache()
	conf := wechat.MiniProgramConfig{Cache: cache}
	app := wechat.NewMiniProgram(auth, conf)
	session, err := app.JsCode2Session(code)
	if err == nil {
		fmt.Println(session.OpenId)
		fmt.Println(session.SessionKey)
		fmt.Println(session.UnionId)
	}
}
