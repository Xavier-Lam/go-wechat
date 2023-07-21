# go-wechat

[![Build Status](https://github.com/Xavier-Lam/go-wechat/actions/workflows/ci.yml/badge.svg)]((https://github.com/Xavier-Lam/go-wechat/actions?query=workflows%3ACI))
[![Coverage Status](https://codecov.io/gh/Xavier-Lam/go-wechat/branch/master/graph/badge.svg)](https://codecov.io/gh/Xavier-Lam/go-wechat)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xavier-Lam/go-wechat)](https://goreportcard.com/report/github.com/Xavier-Lam/go-wechat)

**go-wechat** is a Go package that provides a client for interacting with the WeChat API. It allows you to send API requests to WeChat and handle the response. Currently, it supports only API requests for [Official account](https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Overview.html) and [Mini program](https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/).

[中文版](README.md)

Features:

* Automatically store and update credentials
* Automatically refresh credential and retry after a credential corrupted
* Full unittest coverage
* Easy use and flexible APIs

## Quickstart
* Call encapsulated Apis

        package main

        import (
        	"encoding/json"

            "github.com/Xavier-Lam/go-wechat"
            "github.com/Xavier-Lam/go-wechat/caches"
            "github.com/Xavier-Lam/go-wechat/officialaccount"
        )

        func main() {
            auth := wechat.NewAuth("appId", "appSecret")
            cache := caches.NewDummyCache()
            conf := client.Config{Cache: cache}
	        oa := officialaccount.New(auth, conf)
            jsConfig, err := oa.Js.GetJsConfig("url", officialaccount.JsConfig{})
            data, err := json.Marshal(jsConfig)
        }

* Call api directly

        package main

        import (
            "github.com/Xavier-Lam/go-wechat"
            "github.com/Xavier-Lam/go-wechat/caches"
            "github.com/Xavier-Lam/go-wechat/client"
        )

        func main() {
            auth := wechat.NewAuth("appId", "appSecret")
            cache := caches.NewDummyCache()
            conf := client.Config{Cache: cache}
            w := client.New(auth, conf)
            data := map[string]interface{}{
                "scene": "value1",
                "width": 430,
            }
            resp, err := w.PostJson("/wxa/getwxacodeunlimit", data, true)
        }

* Get latest access token

        w := client.New(auth, conf)
        token, err := w.GetAccessToken()
        ak := token.GetAccessToken()