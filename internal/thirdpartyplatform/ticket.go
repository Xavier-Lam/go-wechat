package thirdpartyplatform

import (
	"errors"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

const (
	BizVerifyTicket              = "verifyticket"
	DefaultVerifyTicketExpiresIn = 43200
)

type (
	VerifyTicketManager        = auth.CredentialManager[string]
	VerifyTicketManagerFactory = func(auth auth.Auth, cache caches.Cache, expiresIn int) VerifyTicketManager
)

type verifyTicketManager struct {
	auth      auth.Auth
	cache     caches.Cache
	expiresIn int
}

func NewVerifyTicketManager(auth auth.Auth, cache caches.Cache, expiresIn int) VerifyTicketManager {
	if expiresIn <= 0 {
		expiresIn = DefaultVerifyTicketExpiresIn
	}
	return &verifyTicketManager{
		auth:      auth,
		cache:     cache,
		expiresIn: expiresIn,
	}
}

func (cm *verifyTicketManager) Get() (*string, error) {
	return cm.get()
}

func (cm *verifyTicketManager) Set(credential *string) error {
	return cm.cache.Set(cm.auth.GetAppId(), BizVerifyTicket, []byte(*credential), cm.expiresIn)
}

func (cm *verifyTicketManager) Renew() (*string, error) {
	return nil, auth.ErrNotRenewable
}

func (cm *verifyTicketManager) Delete() error {
	ticket, err := cm.get()
	if err != nil {
		return err
	}
	return cm.cache.Delete(cm.auth.GetAppId(), BizVerifyTicket, []byte(*ticket))
}

func (m *verifyTicketManager) get() (*string, error) {
	if m.cache == nil {
		return nil, caches.ErrCacheNotSet
	}

	cachedValue, err := m.cache.Get(m.auth.GetAppId(), BizVerifyTicket)
	if err != nil {
		return nil, err
	}

	rv := string(cachedValue)
	if rv == "" {
		return nil, errors.New("empty value")
	}

	return &rv, nil
}
