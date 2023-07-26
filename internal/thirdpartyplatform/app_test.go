package thirdpartyplatform_test

import (
	"testing"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
	"github.com/stretchr/testify/assert"
)

func TestAppSetTicket(t *testing.T) {
	conf := thirdpartyplatform.Config{
		Cache: caches.NewDummyCache(),
	}
	app := thirdpartyplatform.New(test.MockAuth, conf)

	ticket, err := app.GetTicket()
	assert.Error(t, err)
	assert.Equal(t, "", ticket)

	expectedTicket := "test-ticket"
	err = app.SetTicket(expectedTicket)
	assert.NoError(t, err)

	ticket, err = app.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, expectedTicket, ticket)

	auth2 := auth.New("app2", "secret")
	app2 := thirdpartyplatform.New(auth2, conf)

	ticket, err = app2.GetTicket()
	assert.Error(t, err)
	assert.Equal(t, "", ticket)

	expectedTicket2 := "test-ticket2"
	err = app2.SetTicket(expectedTicket2)
	assert.NoError(t, err)

	ticket, err = app2.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, expectedTicket2, ticket)

	ticket, err = app.GetTicket()
	assert.NoError(t, err)
	assert.Equal(t, expectedTicket, ticket)
}
