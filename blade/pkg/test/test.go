package test

import (
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

}

func (h *TestHandler) getDevice(context *gin.Context) {
	device, err := h.webDeviceManager.GetWebDevice()
	routers.ToJson(context, device, err)
}
