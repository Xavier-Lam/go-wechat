package wechat

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount"
)

// Exported interfaces
type (
	Auth              = auth.Auth
	AccessToken       = auth.AccessToken
	AccessTokenClient = client.AccessTokenClient
	WeChatClient      = client.WeChatClient
)

// Exported constructors
var (
	NewAuth              = auth.NewAuth
	NewAccessTokenClient = client.NewAccessTokenClient
	NewMiniProgram       = miniprogram.New
	NewOfficeAccount     = officialaccount.New
	NewWeChatClient      = client.New

	NewAccessToken = auth.NewAccessToken
)

// Exported configurations
type (
	MiniProgramConfig     = miniprogram.Config
	OfficialAccountConfig = officialaccount.Config
	WeChatClientConfig    = client.Config
)

// Exported functions
var (
	GetJson = client.GetJson
)
