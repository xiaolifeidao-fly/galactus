package main

import (
	_ "SignIos/routers"
	"encoding/json"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var AuthKey = ""
var IsOpenAuth = false

func GetIsOpenAuth() {
	IsOpen, err := beego.AppConfig.Bool("auth_open")
	if err != nil {
		IsOpenAuth = true
		return
	}
	IsOpenAuth = IsOpen
}
func GetAuthKey() {
	AuthKey = beego.AppConfig.String("auth_key")
}

type Result struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	ServerTime string      `json:"server_time"`
}

var BeforeExecAllyFunc = func(ctx *context.Context) {
	if IsOpenAuth {
		auth := strings.TrimSpace(ctx.Request.Header.Get("auth"))
		if auth != AuthKey {
			a := Result{
				401,
				"鉴权失败!",
				"",
				time.Now().Format("2006-01-02 15:04:05"),
			}
			message, _ := json.Marshal(a)
			ctx.Output.SetStatus(200)
			ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
			ctx.Output.Body(message)
		}
	}
}

func main() {
	GetIsOpenAuth()
	if IsOpenAuth {
		GetAuthKey()
	}
	beego.InsertFilter("/v1/*", beego.BeforeRouter, BeforeExecAllyFunc)
	beego.Run()
}
