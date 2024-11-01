package device

import (
	"galactus/blade/internal/service/device"
	"galactus/common/base/vo"
	"galactus/common/converter"
	"galactus/common/middleware/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebDevice struct {
	vo.Base
}

func Routers(engine *gin.RouterGroup) {
	//用户列表
	engine.GET("/devices/list", list)

}

func list(context *gin.Context) {
	devices, _ := device.ListWebDevice()
	context.JSON(http.StatusOK, gin.H{
		"code":  routers.SuccessCode,
		"data":  converter.ToVOs[WebDevice](devices),
		"error": nil,
	})
}
