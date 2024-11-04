package dy

import (
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/support/dy"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

func UserRouters(engine *gin.RouterGroup) {
	engine.GET("/dy/user/get", getUserBySecUid)

}

func getUserBySecUid(context *gin.Context) {
	secUid := context.Query("secUid")
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	userInfo, _ := dy.GetUserInfoByWeb(secUid, webDeviceDTO)
	routers.ToJson(context, userInfo, nil)
}
