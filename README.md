# go-wechat

[![Build Status](https://github.com/Xavier-Lam/go-wechat/actions/workflows/ci.yml/badge.svg)]((https://github.com/Xavier-Lam/go-wechat/actions?query=workflows%3ACI))
[![Coverage Status](https://codecov.io/gh/Xavier-Lam/go-wechat/branch/master/graph/badge.svg)](https://codecov.io/gh/Xavier-Lam/go-wechat)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xavier-Lam/go-wechat)](https://goreportcard.com/report/github.com/Xavier-Lam/go-wechat)

**go-wechat** 是一个Go拓展，提供了一个用于与微信API进行交互的客户端。它允许您向微信发送API请求并处理响应。目前支持[公众号](https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Overview.html)和[小程序](https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/).

**[English Readme](README.en.md)**

功能:

* 自动存储和更新Access token
* Access token失效后自动刷新并重试
* 完整的单元测试
* 易用灵活的接口

## 快速开始
* 调用封装API

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

* 调用原始API

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

* 获取最新token

        w := client.New(auth, conf)
        token, err := w.GetAccessToken()
        ak := token.GetAccessToken()