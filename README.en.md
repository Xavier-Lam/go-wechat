# go-wechat

[![Go Report Card](https://goreportcard.com/badge/github.com/Xavier-Lam/go-wechat)](https://goreportcard.com/report/github.com/Xavier-Lam/go-wechat)

**go-wechat** is a Go package that provides a client for interacting with the WeChat API. It allows you to send API requests to WeChat and handle the response. Currently, it supports only API requests for [Official account](https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Overview.html) and [Mini program](https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/).

[中文版](README.md)

Features:

* Automatically store and update credentials
* Automatically refresh credential and retry after a credential corrupted
* Full unittest coverage
* Easy use and flexible APIs

## Quickstart
* Basic usage

        package main

        import (
            "github.com/Xavier-Lam/go-wechat"
            "github.com/Xavier-Lam/go-wechat/caches"
            "github.com/Xavier-Lam/go-wechat/client"
        )

        func main() {
            auth := wechat.NewAuth("appId", "appSecret")
            cache := caches.NewDummyCache()
            conf := &client.Config{Cache: cache}
            c := client.New(auth, conf)
            data := map[string]interface{}{
                "scene": "value1",
                "width": 430,
            }
            resp, err := c.PostJson("/wxa/getwxacodeunlimit", data, true)
        }

* Get latest access token

        c := client.New(auth, conf)
        token, err := c.GetAccessToken()
        ak := token.GetAccessToken()