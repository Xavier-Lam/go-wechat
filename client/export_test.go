package client

import "net/url"

func NewAccessTokenClient(url *url.URL, http HttpClient) AccessTokenClient {
	return &accessTokenClient{
		http: http,
		url:  url,
	}
}
