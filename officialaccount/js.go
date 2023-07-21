package officialaccount

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"

	"github.com/Xavier-Lam/go-wechat"
	"github.com/Xavier-Lam/go-wechat/caches"
	"github.com/Xavier-Lam/go-wechat/officialaccount/apis"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type JsConfig struct {
	Debug     bool     `json:"debug"`
	AppId     string   `json:"appId"`
	Timestamp int      `json:"timestamp"`
	NonceStr  string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

type js struct {
	api   apis.Js
	auth  wechat.Auth
	cache caches.Cache
}

func newJs(auth wechat.Auth, api apis.Js, cache caches.Cache) *js {
	return &js{
		api:   api,
		auth:  auth,
		cache: cache,
	}
}

// Get the latest validate ticket (obtaining from cache first)
// It may return an error along with the ticket if there is no `Cache` set up.
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
func (j *js) GetTicket() (string, error) {
	if j.cache != nil {
		cachedValue, err := j.cache.Get(j.auth.GetAppId(), caches.BizJSTicket)
		if err == nil {
			if ticket := string(cachedValue); ticket != "" {
				return ticket, nil
			}
		}
	}

	return j.FetchTicket()
}

// Obtaining api_ticket from server side
// It may return an error along with the ticket if there is no `Cache` set up.
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#54
func (j *js) FetchTicket() (string, error) {
	ticket, err := j.api.GetTicket()
	if err != nil {
		return "", err
	}

	if j.cache == nil {
		err = fmt.Errorf("cache is not set")
	} else {
		err = j.cache.Set(
			j.auth.GetAppId(),
			caches.BizJSTicket,
			[]byte(ticket.Ticket),
			ticket.ExpiresIn,
		)
	}

	return ticket.Ticket, err
}

func (j *js) GetJsConfig(url string, c JsConfig) (JsConfig, error) {
	var err error
	c.AppId = j.auth.GetAppId()
	if c.NonceStr == "" {
		c.NonceStr = getRandomString(8)
	}
	if c.Timestamp <= 0 {
		c.Timestamp = int(time.Now().Unix())
	}
	if c.JsApiList == nil {
		c.JsApiList = []string{}
	}
	c.Signature, err = j.Sign(url, c.NonceStr, c.Timestamp)
	if err != nil {
		return JsConfig{}, err
	}

	return c, nil
}

func (j *js) Sign(url string, nonceStr string, timestamp int) (string, error) {
	ticket, err := j.GetTicket()
	if err != nil {
		return "", err
	}
	strToSign := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonceStr, timestamp, url,
	)
	hash := sha1.New()
	_, err = hash.Write([]byte(strToSign))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
