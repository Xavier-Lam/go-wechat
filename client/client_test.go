package client_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/client"
	"github.com/stretchr/testify/assert"
)

type RequestHandler func(req *http.Request, calls int) (*http.Response, error)

type mockHttpClient struct {
	calls   int
	handler RequestHandler
}

func newMockHttpClient(handler RequestHandler) client.HttpClient {
	return &mockHttpClient{handler: handler}
}

func (c *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	c.calls++
	return c.handler(req, c.calls)
}

type mockAccessTokenClient struct {
	token string
}

func newMockAccessTokenClient(token string) client.AccessTokenClient {
	return &mockAccessTokenClient{token: token}
}

func (c *mockAccessTokenClient) GetAccessToken(auth wechat.Auth) (client.Token, error) {
	return client.NewToken(c.token, client.DefaultTokenExpiresIn), nil
}

var (
	emptyResponse = httptest.NewRecorder().Result()
)

func createEmptyResponse() (*http.Response, error) {
	return emptyResponse, nil
}

func createJsonResponse(json string) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	recorder.Header().Add("Content-Type", "application/json")
	recorder.WriteString(json)
	return recorder.Result(), nil
}

var (
	appID     = "mock-app-id"
	appSecret = "mock-app-secret"
	auth      = wechat.NewAuth(appID, appSecret)
)

func TestWeChatClientGet(t *testing.T) {
	url := "https://api.weixin.qq.com/some-endpoint?a=1"
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, url, req.URL.String())

		return createEmptyResponse()
	})

	config := &client.Config{HttpClient: mc}
	c := client.New(auth, config)

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
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
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

		return createEmptyResponse()
	})

	config := &client.Config{HttpClient: mc}
	c := client.New(auth, config)

	resp, err := c.PostJson(url, data, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientDo(t *testing.T) {
	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.URL.String(), req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return createEmptyResponse()
	})

	config := &client.Config{HttpClient: mc}
	c := client.New(auth, config)

	resp, err := c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	// relative url
	expectedBaseUrl, _ := url.Parse("https://example.com")
	expectedFullUrl := "https://example.com/some-endpoint"
	expectedRelativeUrl := "/some-endpoint"
	expectedReq, _ = http.NewRequest("GET", expectedRelativeUrl, nil)
	mc = newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedFullUrl, req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return createEmptyResponse()
	})

	config = &client.Config{
		HttpClient: mc,
		BaseApiUri: expectedBaseUrl,
	}
	c = client.New(auth, config)

	resp, err = c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	// default base url, relative url
	expectedReq, _ = http.NewRequest("GET", expectedUrl, nil)
	mc = newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedUrl, req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)

		return createEmptyResponse()
	})

	config = &client.Config{
		HttpClient: mc,
		BaseApiUri: expectedBaseUrl,
	}
	c = client.New(auth, config)

	resp, err = c.Do(expectedReq, false)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientWithToken(t *testing.T) {
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.URL.String(), req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)
		assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

		return createEmptyResponse()
	})

	cache := caches.NewDummyCache()
	cache.Set(appID, caches.CacheBizAccessToken, client.NewToken(accessToken, 3600), 3600)
	config := &client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(auth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)
}

func TestWeChatClientDoWithoutToken(t *testing.T) {
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		assert.Equal(t, 1, calls)
		assert.Equal(t, expectedReq.URL.String(), req.URL.String())
		assert.Equal(t, expectedReq.Method, req.Method)
		assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

		return createEmptyResponse()
	})

	akc := newMockAccessTokenClient(accessToken)
	cache := caches.NewDummyCache()
	config := &client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.NewWithDependency(akc, auth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	storedToken, err := cache.Get(appID, caches.CacheBizAccessToken)
	assert.NoError(t, err)
	assert.Equal(t, storedToken.(client.Token).GetAccessToken(), accessToken)
}

func TestWeChatClientDoWithInvalidToken(t *testing.T) {
	invalidToken := "invalid"
	accessToken := "token"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		if calls == 1 {
			assert.Equal(t, expectedReq.URL.String(), req.URL.String())
			assert.Equal(t, expectedReq.Method, req.Method)
			assert.Equal(t, invalidToken, req.URL.Query().Get("access_token"))

			return createJsonResponse(`{"errcode": 40014, "errmsg": "Invalid access token"}`)
		} else if calls == 2 {
			assert.Equal(t, expectedReq.URL.String(), req.URL.String())
			assert.Equal(t, expectedReq.Method, req.Method)
			assert.Equal(t, accessToken, req.URL.Query().Get("access_token"))

			return createEmptyResponse()
		} else {
			assert.Fail(t, "Unexpected calls")
			return nil, nil
		}
	})

	akc := newMockAccessTokenClient(accessToken)
	cache := caches.NewDummyCache()
	cache.Set(appID, caches.CacheBizAccessToken, client.NewToken(invalidToken, 3600), 3600)
	config := &client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.NewWithDependency(akc, auth, config)

	resp, err := c.Do(expectedReq, true)
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, resp)

	storedToken, err := cache.Get(appID, caches.CacheBizAccessToken)
	assert.NoError(t, err)
	assert.Equal(t, storedToken.(client.Token).GetAccessToken(), accessToken)
}

func TestWeChatClientDoWithInvalidTokenAndInvalidCredential(t *testing.T) {
	invalidToken := "invalid"

	expectedUrl := "https://api.weixin.qq.com/some-endpoint"
	expectedReq, _ := http.NewRequest("GET", expectedUrl, nil)
	mc := newMockHttpClient(func(req *http.Request, calls int) (*http.Response, error) {
		if calls == 1 {
			assert.Equal(t, expectedReq.URL.String(), req.URL.String())
			assert.Equal(t, expectedReq.Method, req.Method)
			assert.Equal(t, invalidToken, req.URL.Query().Get("access_token"))

			return createJsonResponse(`{"errcode": 40014, "errmsg": "Invalid access token"}`)
		} else if calls == 2 {
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "https", req.URL.Scheme)
			assert.Equal(t, "api.weixin.qq.com", req.URL.Host)
			assert.Equal(t, "/cgi-bin/token", req.URL.Path)
			assert.Equal(t, "client_credential", req.URL.Query().Get("grant_type"))
			assert.Equal(t, appID, req.URL.Query().Get("appid"))
			assert.Equal(t, appSecret, req.URL.Query().Get("secret"))

			return createJsonResponse(`{"errcode": 40125, "errmsg": "invalid appsecret"}`)
		} else {
			assert.Fail(t, "Unexpected calls")
			return nil, nil
		}
	})

	cache := caches.NewDummyCache()
	cache.Set(appID, caches.CacheBizAccessToken, client.NewToken(invalidToken, 3600), 3600)
	config := &client.Config{
		HttpClient: mc,
		Cache:      cache,
	}
	c := client.New(auth, config)

	_, err := c.Do(expectedReq, true)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "40014")
	assert.ErrorContains(t, err, "40125")
	assert.Equal(t, err.(client.WeChatApiError).ErrCode, 40014)
	assert.Equal(t, err.(client.WeChatApiError).RetryError.(client.WeChatApiError).ErrCode, 40125)
}

func TestWeChatClientGetAccessToken(t *testing.T) {
	oldToken := "old"
	newToken := "token"

	cache := caches.NewDummyCache()
	config := &client.Config{Cache: cache}
	akc := newMockAccessTokenClient(oldToken)
	c := client.NewWithDependency(akc, auth, config)

	token, err := c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), oldToken)

	token, err = c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), oldToken)

	akc = newMockAccessTokenClient(newToken)
	c = client.NewWithDependency(akc, auth, config)

	token, err = c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), oldToken)

	token, err = c.FetchAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), newToken)

	token, err = c.GetAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, token.GetAccessToken(), newToken)
}

func TestWeChatClientGetAppId(t *testing.T) {
	client := client.New(auth, nil)

	result := client.GetAuth()
	assert.Equal(t, appID, result.GetAppId())
	assert.Equal(t, appSecret, result.GetAppSecret())
}
