package test

import (
	"galactus/blade/internal/consts"
	"galactus/blade/internal/service/device/biz"
	biz2 "galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	*routers.BaseHandler
	webDeviceManager *biz.WebDeviceManager
	ipManager        *biz2.IpManager
}

func NewTestHandler() *TestHandler {
	return &TestHandler{
		webDeviceManager: biz.GetDefaultWebDeviceManager(),
		ipManager:        biz2.GetDefaultIpManager(),
	}
}

func (h *TestHandler) RegisterHandler(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/test/device", h.getDevice)
	//获取ip
	engine.GET("/test/ip", h.getIp)
}

func (h *TestHandler) getDevice(context *gin.Context) {
	//device, err := h.webDeviceManager.GetWebDevice()
	device, err := biz.GetDefaultWebDeviceManager().GetWebDevice(consts.SceneAuditLike)
	routers.ToJson(context, device, err)
}

func (h *TestHandler) getIp(context *gin.Context) {
	ip, err := biz2.GetDefaultIpManager().GetIp(consts.SceneCollectDevice)
	routers.ToJson(context, ip, err)
}
