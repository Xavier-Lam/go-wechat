package auth

// Auth represents the authentication interface
type Auth interface {
	// AppId of your WeChat application
	GetAppId() string
	// AppSecret of your WeChat application
	GetAppSecret() string
}

// auth implements the Auth interface
type auth struct {
	appId     string
	appSecret string
}

// NewAuth creates a new instance of Auth
func NewAuth(appId string, appSecret string) Auth {
	return &auth{
		appId:     appId,
		appSecret: appSecret,
	}
}

// GetAppID returns the AppId
func (a *auth) GetAppId() string {
	return a.appId
}

// GetAppSecret returns the AppSecret
func (a *auth) GetAppSecret() string {
	return a.appSecret
}
