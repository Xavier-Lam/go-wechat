package thirdpartyplatform_test

import (
	"errors"
	"testing"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
	"github.com/stretchr/testify/assert"
)

const (
	mockTicket = "mock-ticket"
)

type mockVerifyTicketManager struct {
	ticket string
}

func (cm *mockVerifyTicketManager) Get() (*string, error) {
	return &cm.ticket, nil
}

func (cm *mockVerifyTicketManager) Set(credential *string) error {
	return errors.New("not implemented")
}

func (cm *mockVerifyTicketManager) Renew() (*string, error) {
	return nil, auth.ErrNotRenewable
}

func (cm *mockVerifyTicketManager) Delete() error {
	return auth.ErrNotDeletable
}

func TestVerifyTicketManager(t *testing.T) {
	var empty *string

	cache := caches.NewDummyCache()
	manager := thirdpartyplatform.NewVerifyTicketManager(test.MockAuth, cache, 3600)

	ticket, err := manager.Get()
	assert.Error(t, err)
	assert.Equal(t, empty, ticket)

	err = manager.Delete()
	assert.Error(t, err)

	expectedTicket := "test-ticket"
	err = manager.Set(&expectedTicket)
	assert.NoError(t, err)

	ticket, err = manager.Get()
	assert.NoError(t, err)
	assert.Equal(t, &expectedTicket, ticket)

	expectedTicket = "test-ticket2"
	err = manager.Set(&expectedTicket)
	assert.NoError(t, err)

	ticket, err = manager.Get()
	assert.NoError(t, err)
	assert.Equal(t, &expectedTicket, ticket)

	manager = thirdpartyplatform.NewVerifyTicketManager(test.MockAuth, cache, 3600)
	ticket, err = manager.Get()
	assert.NoError(t, err)
	assert.Equal(t, &expectedTicket, ticket)

	err = manager.Delete()
	assert.NoError(t, err)

	ticket, err = manager.Get()
	assert.Error(t, err)
	assert.Equal(t, empty, ticket)

	ticket, err = manager.Renew()
	assert.ErrorIs(t, err, auth.ErrNotRenewable)
	assert.Equal(t, empty, ticket)
}
