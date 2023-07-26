package thirdpartyplatform_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestPrepareRequest(t *testing.T) {
	c := thirdpartyplatform.NewAccessTokenClient(nil, "")
	vtm := mockVerifyTicketManager{mockTicket}
	a := thirdpartyplatform.NewThirdPartyPlatformAuth(test.MockAuth, &vtm)

	req, err := c.PrepareRequest(a)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	test.AssertEndpointEqual(t, thirdpartyplatform.DefaultAccessTokenUrl, req.URL)

	body := verifyTicketRequest{}
	err = json.NewDecoder(req.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, test.AppId, body.AppId)
	assert.Equal(t, test.AppSecret, body.AppSecret)
	assert.Equal(t, mockTicket, body.VerifyTicket)

	requestUrl := "https://example.com/test"
	c = thirdpartyplatform.NewAccessTokenClient(nil, requestUrl)

	req, err = c.PrepareRequest(a)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	test.AssertEndpointEqual(t, requestUrl, req.URL)

	body = verifyTicketRequest{}
	err = json.NewDecoder(req.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, test.AppId, body.AppId)
	assert.Equal(t, test.AppSecret, body.AppSecret)
	assert.Equal(t, mockTicket, body.VerifyTicket)
}

func TestSendRequest(t *testing.T) {
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

		return test.Responses.Empty()
	})
	c := thirdpartyplatform.NewAccessTokenClient(httpClient, "")
	vtm := mockVerifyTicketManager{mockTicket}
	a := thirdpartyplatform.NewThirdPartyPlatformAuth(test.MockAuth, &vtm)

	req, err := c.PrepareRequest(a)
	assert.NoError(t, err)
	resp, err := c.SendRequest(a, req)

	assert.NoError(t, err)
	assert.Equal(t, httptest.NewRecorder().Result(), resp)
}

func TestHandleResponse(t *testing.T) {
	c := thirdpartyplatform.NewAccessTokenClient(nil, "")
	vtm := mockVerifyTicketManager{mockTicket}
	a := thirdpartyplatform.NewThirdPartyPlatformAuth(test.MockAuth, &vtm)
	resp, _ := test.Responses.Json(`{"component_access_token": "access-token", "expires_in": 7200}`)

	token, err := c.HandleResponse(a, resp, &http.Request{})

	assert.NoError(t, err)
	assert.Equal(t, "access-token", token.GetAccessToken())
	assert.Equal(t, 7200, token.GetExpiresIn())
}
