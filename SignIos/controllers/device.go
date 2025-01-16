package controllers

import (
	"SignIos/common"
	"encoding/base64"
	"encoding/json"
)

type DeviceController struct {
	BaseController
}

type DeviceReq struct {
	Base64 string `json:"base_64"`
}

func (c *DeviceController) EncodeData() {
	var v DeviceReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Base64 == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	data, err := base64.StdEncoding.DecodeString(v.Base64)
	if err != nil {
		c.responseFail("解码base64失败!", nil)
	}
	err, str := common.EncryptToBase64(data)
	if err != nil {
		c.responseFail("加密失败!", nil)
	}
	c.responseSuccess("加密成功!", str)
}

func (c *DeviceController) DecodeData() {
	var v DeviceReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Base64 == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	err, str := common.Base64StrDecrypt(v.Base64)
	if err != nil {
		c.responseFail("解密失败,请查看base64加密方式是否正确,请使用url模式base64!", nil)
	}
	c.responseSuccess("解密成功!", str)
}
