package wechat

import (
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/miniprogram"
	"github.com/Xavier-Lam/go-wechat/officialaccount"
)

var (
	NewAuth          = auth.NewAuth
	NewMiniProgram   = miniprogram.New
	NewOfficeAccount = officialaccount.New
	NewWeChatClient  = client.New
)

type (
	MiniProgramConfig     = miniprogram.Config
	OfficialAccountConfig = officialaccount.Config
	WeChatClientConfig    = client.Config
)
