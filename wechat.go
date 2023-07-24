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
	AccessTokenClient = client.AccessTokenClient
	WeChatClient      = client.WeChatClient
)

// Exported factories
var (
	AccessTokenClientFactory  = client.AccessTokenClientFactory
	AccessTokenManagerFactory = client.AccessTokenManagerFactory
)

// Exported constructors
var (
	NewAuth          = auth.New
	NewMiniProgram   = miniprogram.New
	NewOfficeAccount = officialaccount.New
	NewWeChatClient  = client.New

	// less commonly used
	NewAccessToken        = auth.NewAccessToken
	NewAccessTokenClient  = client.NewAccessTokenClient
	NewAccessTokenManager = client.NewAccessTokenManager
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
