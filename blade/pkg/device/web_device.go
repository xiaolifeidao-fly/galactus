package device

import (
	deviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/common/converter"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

type WebDeviceHandler struct {
	*routers.BaseHandler
	webDeviceService *deviceService.WebDeviceService
}

func NewWebDeviceHandler() *WebDeviceHandler {
	return &WebDeviceHandler{
		webDeviceService: deviceService.NewWebDeviceService(),
	}
}

func (h *WebDeviceHandler) RegisterHandler(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/devices/list", h.list)
	engine.POST("/devices/save", h.save)

}

func (h *WebDeviceHandler) save(context *gin.Context) {
	var device dto.WebDeviceDTO
	context.ShouldBindJSON(&device)
	deviceDTO := converter.ToDTO[dto.WebDeviceDTO](&device)
	_, err := h.webDeviceService.Save(deviceDTO)
	routers.ToJson(context, "保存成功", err)
}

func (h *WebDeviceHandler) list(context *gin.Context) {
	devices, err := h.webDeviceService.GetWebDeviceList()
	routers.ToJson(context, devices, err)
}
