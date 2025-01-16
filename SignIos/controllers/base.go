package controllers

import (
	"github.com/astaxie/beego"
	"time"
)

type BaseController struct {
	beego.Controller
}

// 请求返回值自定义方法
type Result struct {
	Code       int `json:"code"`
	Message    string `json:"message"`
	Data       interface{} `json:"data"`
	ServerTime string `json:"server_time"`
}

const (
	SUCCESS = iota
	FAIL
)

func SendResponse(code int, message string, data interface{}) (result Result) {
	result.Code = code
	result.Message = message
	result.Data = data
	result.ServerTime = time.Now().Format("2006-01-02 15:04:05")
	return
}
func (c *BaseController) responseJson(code int, msg string, data interface{}) {
	result := map[string]interface{}{"code": code, "msg": msg, "data": data, "server_time": time.Now().Format("2006-01-02 15:04:05")}
	c.Data["json"] = result
	c.ServeJSON()
	c.StopRun()
}
func (c *BaseController) responseSuccess(msg string, data interface{}) {
	c.responseJson(1, msg, data)
}
func (c *BaseController) responseFail(msg string, data interface{}) {
	c.responseJson(0, msg, data)
}

