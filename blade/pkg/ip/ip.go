package ip

import (
	"galactus/blade/internal/consts"
	ip "galactus/blade/internal/service/ip/biz"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

type IpHandler struct {
	*routers.BaseHandler
	IpManager *ip.IpManager
}

func NewIpHandler() *IpHandler {
	return &IpHandler{
		IpManager: ip.GetDefaultIpManager(),
	}
}

func (h *IpHandler) RegisterHandler(engine *gin.RouterGroup) {
	engine.GET("/ip/collectDeviceIp", h.getCollectDeviceIp)
}

func (h *IpHandler) getCollectDeviceIp(context *gin.Context) {
	ip, err := h.IpManager.GetIp(consts.SceneCollectDevice)
	if err != nil {
		routers.ToJson(context, nil, err)
		return
	}
	routers.ToJson(context, ip.Ip, nil)
}
