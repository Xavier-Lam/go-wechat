package auth

// Auth represents the authentication interface
type Auth interface {
	// AppId of your WeChat application
	GetAppId() string
	// AppSecret of your WeChat application
	GetAppSecret() string
}

// wechatAuth implements the Auth interface
type wechatAuth struct {
	appId     string
	appSecret string
}

// NewAuth creates a new instance of Auth
func NewAuth(appId string, appSecret string) Auth {
	return &wechatAuth{
		appId:     appId,
		appSecret: appSecret,
	}
}

// GetAppID returns the AppId
func (a *wechatAuth) GetAppId() string {
	return a.appId
}

// GetAppSecret returns the AppSecret
func (a *wechatAuth) GetAppSecret() string {
	return a.appSecret
}
