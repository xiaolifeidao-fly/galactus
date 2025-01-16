// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"SignIos/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/dy", &controllers.DyController{}, "post:GetParams"),
		beego.NSRouter("/dy_decode", &controllers.DyController{}, "post:DecryptInfo"),
		beego.NSRouter("/jisu", &controllers.DyJisuController{}, "post:GetParams"),
		beego.NSRouter("/jisu_decode", &controllers.DyJisuController{}, "post:DecryptInfo"),
		beego.NSRouter("/device_encode", &controllers.DeviceController{}, "post:EncodeData"),
		beego.NSRouter("/device_decode", &controllers.DeviceController{}, "post:DecodeData"),
		beego.NSRouter("/mssdk_token_req_decode", &controllers.MssdkController{}, "post:DecodeTokenReqData"),
		beego.NSRouter("/mssdk_token_res_decode", &controllers.MssdkController{}, "post:DecodeTokenResData"),
		beego.NSRouter("/jisu_mssdk_token_req_decode", &controllers.MssdkController{}, "post:JiSuDecodeTokenReqData"),
		beego.NSRouter("/jisu_mssdk_token_res_decode", &controllers.MssdkController{}, "post:JiSuDecodeTokenResData"),
		beego.NSRouter("/duoshan", &controllers.DuoShanController{}, "post:GetParams"),
		beego.NSRouter("/duoshan_decode", &controllers.DuoShanController{}, "post:DecryptInfo"),
	)
	beego.AddNamespace(ns)
}
