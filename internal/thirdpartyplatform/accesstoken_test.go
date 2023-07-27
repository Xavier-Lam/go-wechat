package thirdpartyplatform_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/Xavier-Lam/go-wechat/internal/thirdpartyplatform"
	"github.com/stretchr/testify/assert"
)

type verifyTicketRequest struct {
	AppId        string `json:"component_appid"`
	AppSecret    string `json:"component_appsecret"`
	VerifyTicket string `json:"component_verify_ticket"`
}

func TestAccessTokenFetcher(t *testing.T) {
	httpClient := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, http.MethodPost, req.Method)
		test.AssertEndpointEqual(t, thirdpartyplatform.DefaultAccessTokenUrl, req.URL)

		body := verifyTicketRequest{}
		err := json.NewDecoder(req.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, test.AppId, body.AppId)
		assert.Equal(t, test.AppSecret, body.AppSecret)
		assert.Equal(t, mockTicket, body.VerifyTicket)

		return test.Responses.Json(`{"component_access_token": "access-token", "expires_in": 7200}`)
	})

	token, err := thirdpartyplatform.ExportedAccessTokenFetcher(httpClient, test.MockAuth, mockTicket, nil)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())

	requestUrl, _ := url.Parse("https://example.com/test")
	httpClient = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, http.MethodPost, req.Method)
		test.AssertEndpointEqual(t, "https://example.com/test", req.URL)

		body := verifyTicketRequest{}
		err := json.NewDecoder(req.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, test.AppId, body.AppId)
		assert.Equal(t, test.AppSecret, body.AppSecret)
		assert.Equal(t, mockTicket, body.VerifyTicket)

		return test.Responses.Json(`{"component_access_token": "access-token", "expires_in": 7200}`)
	})

	token, err = thirdpartyplatform.ExportedAccessTokenFetcher(httpClient, test.MockAuth, mockTicket, requestUrl)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())

	httpClient = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		return test.Responses.Empty()
	})

	token, err = thirdpartyplatform.ExportedAccessTokenFetcher(httpClient, test.MockAuth, mockTicket, nil)
	assert.Error(t, err)
	assert.Nil(t, token)
}
