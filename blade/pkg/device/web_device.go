package device

import (
	deviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/dto"
	"galactus/common/converter"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

func WebDeviceRouters(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/devices/list", list)
	engine.POST("/devices/save", save)

}

func save(context *gin.Context) {
	var device dto.WebDeviceDTO
	context.ShouldBindJSON(&device)
	deviceDTO := converter.ToDTO[dto.WebDeviceDTO](&device)
	_, err := deviceService.NewWebDeviceService().Save(deviceDTO)
	routers.ToJson(context, "保存成功", err)
}

func list(context *gin.Context) {
	devices, err := deviceService.NewWebDeviceService().GetWebDeviceList()
	routers.ToJson(context, devices, err)
}
