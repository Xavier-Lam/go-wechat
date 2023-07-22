package main

import (
	"fmt"
	"os"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
)

func main() {
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
		fmt.Println(session)
	}
}
