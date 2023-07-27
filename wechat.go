package wechat

import (
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/miniprogram"
	"github.com/Xavier-Lam/go-wechat/internal/officialaccount"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
)

// Exported interfaces
type (
	Auth                = auth.Auth
	VeryfyTicketManager = thirdpartyplatform.VerifyTicketManager
	WeChatClient        = client.WeChatClient
)

type (
	AccessToken = auth.AccessToken
)

// Exported constructors
var (
	NewAuth               = auth.New
	NewMiniProgram        = miniprogram.New
	NewOfficeAccount      = officialaccount.New
	NewThirdPartyPlatform = thirdpartyplatform.New
	NewWeChatClient       = client.New

	// less commonly used
	NewAccessToken = auth.NewAccessToken
)

// Exported configurations
type (
	MiniProgramConfig        = miniprogram.Config
	OfficialAccountConfig    = officialaccount.Config
	ThirdPartyPlatformConfig = thirdpartyplatform.Config
	WeChatClientConfig       = client.Config
)

// Exported functions
var (
	GetJson = client.GetJson
)
