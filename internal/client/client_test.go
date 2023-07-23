package client_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/internal/auth"
	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

var (
	emptyResponse = httptest.NewRecorder().Result()
)

var (
	appID     = "mock-app-id"
	appSecret = "mock-app-secret"
	mockAuth  = auth.NewAuth(appID, appSecret)
)

func TestWeChatClientGet(t *testing.T) {
	url := "https://api.weixin.qq.com/some-endpoint?a=1"
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, url, req.URL.String())

		return test.Responses.Empty()
	})

	config := client.Config{HttpClient: mc}
	c := client.New(mockAuth, config)

	resp, err := c.Get(url, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientPost(t *testing.T) {
	url := "https://api.weixin.qq.com/some-endpoint?a=1"
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, url, req.URL.String())
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		bodyBytes, _ := json.Marshal(data)

		expectedBody := bytes.NewReader(bodyBytes)
		actualBody := req.Body

		buffer := make([]byte, len(bodyBytes))
		_, err := actualBody.Read(buffer)
		if err != nil {
			t.Errorf("Error reading request body: %v", err)
			return nil, err
		}
		actualBody.Close()

		assert.Equal(t, expectedBody, bytes.NewReader(buffer))

		return test.Responses.Empty()
	})

	config := client.Config{HttpClient: mc}
	c := client.New(mockAuth, config)

	resp, err := c.PostJson(url, data, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientDo(t *testing.T) {
	expectedUrl := "https://api.weixin.qq.com/some-endpoint?a=1"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.URL.String(), req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return test.Responses.Empty()
	})

	config := client.Config{HttpClient: mc}
	c := client.New(mockAuth, config)

	resp, err := c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	// relative url
	expectedBaseUrl, _ := url.Parse("https://example.com")
	expectedFullUrl := "https://example.com/some-endpoint"
	expectedRelativeUrl := "/some-endpoint"
	expectedReq, _ = http.NewRequest("GET", expectedRelativeUrl, nil)
	mc = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedFullUrl, req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return test.Responses.Empty()
	})

	config = client.Config{
		HttpClient: mc,
		BaseApiUrl: expectedBaseUrl,
	}
	c = client.New(mockAuth, config)

	resp, err = c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	// default base url, relative url
	expectedReq, _ = http.NewRequest("GET", expectedUrl, nil)
	mc = test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedUrl, req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return test.Responses.Empty()
	})

	config = client.Config{
		HttpClient: mc,
		BaseApiUrl: expectedBaseUrl,
	}
	c = client.New(mockAuth, config)

	resp, err = c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientWithToken(t *testing.T) {
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.Method, req.Method)
		test.AssertEndpointEqual(t, expectedUrl, req.URL)
		assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

		return test.Responses.Empty()
	})

	cache := caches.NewDummyCache()
	serializedToken, _ := auth.SerializeToken(auth.NewAccessToken(accessToken, 3600))
	cache.Set(appID, caches.BizAccessToken, serializedToken, 3600)
	config := client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(mockAuth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientDoWithoutToken(t *testing.T) {
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.Method, req.Method)
		test.AssertEndpointEqual(t, expectedUrl, req.URL)
		assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

		return test.Responses.Empty()
	})

	atc := test.NewMockAccessTokenClient(accessToken)
	cache := caches.NewDummyCache()
	config := client.Config{
		CredentialManagerFactory: func(auth auth.Auth, c http.Client, cache caches.Cache, accessTokenUrl *url.URL) auth.CredentialManager {
			return client.NewWeChatAccessTokenCredentialManager(auth, cache, atc)
		},
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(mockAuth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	storedToken, err := cache.Get(appID, caches.BizAccessToken)
	assert.NoError(t, err)
	token, err := auth.DeserializeToken(storedToken)
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), accessToken)
}

func TestWeChatClientDoWithInvalidToken(t *testing.T) {
	invalidToken := "invalid"
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		if calls == 1 {
			assert.Equal(t, expectedReq.Method, req.Method)
			test.AssertEndpointEqual(t, expectedUrl, req.URL)
			assert.Equal(t, invalidToken, req.URL.Query().Get("access_token"))

			return test.Responses.Json(`{"errcode": 40014, "errmsg": "Invalid access token"}`)
		} else if calls == 2 {
			assert.Equal(t, expectedReq.Method, req.Method)
			test.AssertEndpointEqual(t, expectedUrl, req.URL)
			assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

			return test.Responses.Empty()
		} else {
			assert.Fail(t, "Unexpected calls")
			return nil, nil
		}
	})

	atc := test.NewMockAccessTokenClient(accessToken)
	cache := caches.NewDummyCache()
	serializedToken, _ := auth.SerializeToken(auth.NewAccessToken(invalidToken, 3600))
	cache.Set(appID, caches.BizAccessToken, serializedToken, 3600)
	config := client.Config{
		CredentialManagerFactory: func(auth auth.Auth, c http.Client, cache caches.Cache, accessTokenUrl *url.URL) auth.CredentialManager {
			return client.NewWeChatAccessTokenCredentialManager(auth, cache, atc)
		},
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(mockAuth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	storedToken, err := cache.Get(appID, caches.BizAccessToken)
	assert.NoError(t, err)
	token, err := auth.DeserializeToken(storedToken)
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), accessToken)
}

func TestWeChatClientDoWithInvalidTokenAndInvalidCredential(t *testing.T) {
	invalidToken := "invalid"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := test.NewMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		if calls == 1 {
			assert.Equal(t, expectedReq.Method, req.Method)
			test.AssertEndpointEqual(t, expectedUrl, req.URL)
			assert.Equal(t, invalidToken, req.URL.Query().Get("access_token"))

			return test.Responses.Json(`{"errcode": 40014, "errmsg": "Invalid access token"}`)
		} else if calls == 2 {
			assert.Equal(t, "GET", req.Method)
			test.AssertEndpointEqual(t, client.DefaultAccessTokenUri, req.URL)
			assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
			assert.Equal(t, appID, req.URL.Query().Get("appid"))
			assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

			return test.Responses.Json(`{"errcode": 40125, "errmsg": "invalid appsecret"}`)
		} else {
			assert.Fail(t, "Unexpected calls")
			return nil, nil
		}
	})

	cache := caches.NewDummyCache()
	serializedToken, _ := auth.SerializeToken(auth.NewAccessToken(invalidToken, 3600))
	cache.Set(appID, caches.BizAccessToken, serializedToken, 3600)
	config := client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(mockAuth, config)

	_, err := c.Do(expectedReq, true)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "40014")
	assert.ErrorContains(t, err, "40125")
	assert.Equal(t, err.(*url.Error).Err.(client.WeChatApiError).ErrCode, 40014)
	assert.Equal(t, err.(*url.Error).Err.(client.WeChatApiError).RetryError.(*url.Error).Err.(client.WeChatApiError).ErrCode, 40125)
}

func TestWeChatClientGetAccessToken(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	atc := test.NewMockAccessTokenClient(oldToken)
	config := client.Config{
		CredentialManagerFactory: func(auth auth.Auth, c http.Client, cache caches.Cache, accessTokenUrl *url.URL) auth.CredentialManager {
			return client.NewWeChatAccessTokenCredentialManager(auth, cache, atc)
		},
		Cache: cache,
	}
	c := client.New(mockAuth, config)

	token, err := c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, oldToken, token.GetAccessToken())

	atc = test.NewMockAccessTokenClient(newToken)
	config = client.Config{
		CredentialManagerFactory: func(auth auth.Auth, c http.Client, cache caches.Cache, accessTokenUrl *url.URL) auth.CredentialManager {
			return client.NewWeChatAccessTokenCredentialManager(auth, cache, atc)
		},
		Cache: cache,
	}
	c = client.New(mockAuth, config)

	token, err = c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, oldToken, token.GetAccessToken())

	token, err = c.FetchAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, newToken, token.GetAccessToken())

	token, err = c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, newToken, token.GetAccessToken())
}

func TestWeChatClientGetAppId(t *testing.T) {
	client := client.New(mockAuth, client.Config{})

	result := client.GetAuth()
	assert.Equal(t, appID, result.GetAppId())
	assert.Equal(t, appSecret, result.GetAppSecret())
}
