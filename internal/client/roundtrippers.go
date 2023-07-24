package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Xavier-Lam/go-wechat/internal/auth"
)

type shouldRetry struct {
	err WeChatApiError
}

func (err shouldRetry) Error() string {
	return err.err.Error()
}

type commonRoundTripper struct {
	baseUrl *url.URL
	next    http.RoundTripper
}

// A common round tripper for every request
func NewCommonRoundTripper(baseUrl *url.URL, next http.RoundTripper) http.RoundTripper {
	return &commonRoundTripper{
		baseUrl: baseUrl,
		next:    next,
	}
}

func (rt *commonRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extend relative url
	if !req.URL.IsAbs() {
		req.URL = rt.baseUrl.ResolveReference(req.URL)
	}

	resp, err := rt.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Handle exceptions
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var apiError WeChatApiError
		err := GetJson(resp, &apiError)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		} else if apiError.ErrCode != 0 {
			return nil, apiError
		}
	}

	return resp, nil
}

type credentialRoundTripper struct {
	cm   auth.CredentialManager
	next http.RoundTripper
}

// A round tripper for credential required requests
func NewCredentialRoundTripper(cm auth.CredentialManager, next http.RoundTripper) http.RoundTripper {
	return &credentialRoundTripper{
		cm:   cm,
		next: next,
	}
}

func (rt *credentialRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req = rt.clearContext(req)
	credentialRequired := req.Context().Value(RequestContextWithCredential) == true

	if credentialRequired {
		req, err = rt.setUpCredential(req, nil)
		if err != nil {
			return
		}
	}

	resp, err = rt.next.RoundTrip(req)

	if credentialRequired && err != nil {
		return rt.handleError(err, req)
	}

	return
}

func (rt *credentialRoundTripper) clearContext(req *http.Request) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, RequestContextCredential, nil)
	return req.WithContext(ctx)
}

func (rt *credentialRoundTripper) setUpCredential(req *http.Request, credential interface{}) (*http.Request, error) {
	var err error
	if credential == nil {
		credential, err = rt.cm.Get()
		if credential == nil {
			return nil, err
		}
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, RequestContextCredential, credential)
	req = req.WithContext(ctx)

	return req, nil
}

func (rt *credentialRoundTripper) handleError(err error, req *http.Request) (*http.Response, error) {
	shouldRetry, ok := err.(shouldRetry)
	if !ok {
		return nil, err
	}

	var credential interface{}
	apiError := shouldRetry.err
	credential, apiError.RetryError = rt.cm.Renew()
	if credential != nil {
		return rt.retry(req, credential)
	}

	return nil, apiError
}

func (rt *credentialRoundTripper) retry(req *http.Request, credential interface{}) (*http.Response, error) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, RequestContextCredential, credential)
	req = req.WithContext(ctx)
	return rt.next.RoundTrip(req)
}

type accessTokenRoundTripper struct {
	next http.RoundTripper
}

// A round tripper for those requests with an access token
func NewAccessTokenRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &accessTokenRoundTripper{
		next: next,
	}
}

func (rt *accessTokenRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	credentialRequired := req.Context().Value(RequestContextWithCredential) == true
	if credentialRequired {
		tokenValue := req.Context().Value(RequestContextCredential)
		token, ok := tokenValue.(*auth.AccessToken)
		if !ok {
			return nil, errors.New("corrupted token")
		}

		query := req.URL.Query()
		query.Set("access_token", token.GetAccessToken())
		req.URL.RawQuery = query.Encode()
	}

	resp, err = rt.next.RoundTrip(req)

	if credentialRequired && err != nil {
		return rt.handleError(err, req)
	}

	return resp, err
}

func (rt *accessTokenRoundTripper) handleError(err error, req *http.Request) (*http.Response, error) {
	apiError, ok := err.(WeChatApiError)
	if !ok {
		return nil, err
	}

	switch apiError.ErrCode {
	case
		ErrCodeAccessTokenExpired,
		ErrCodeInvalidAccessToken,
		ErrCodeInvalidCredential:
		return nil, shouldRetry{err: apiError}
	}

	return nil, err
}

type fetchAccessTokenRoundTripper struct {
	next http.RoundTripper
}

// A round tripper to get the latest access token
func NewFetchAccessTokenRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &fetchAccessTokenRoundTripper{
		next: next,
	}
}

func (rt *fetchAccessTokenRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	credentialRequired := req.Context().Value(RequestContextWithCredential) == true

	if credentialRequired {
		tokenValue := req.Context().Value(RequestContextCredential)
		auth, ok := tokenValue.(auth.Auth)
		if !ok {
			return nil, errors.New("corrupted auth")
		}

		query := url.Values{
			"grant_type": {"client_credential"},
			"appid":      {auth.GetAppId()},
			"secret":     {auth.GetAppSecret()},
		}
		req.URL.RawQuery = query.Encode()
	}

	return rt.next.RoundTrip(req)
}
