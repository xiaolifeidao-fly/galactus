package routers

import (
	"galactus/blade/pkg/device"
	"galactus/blade/pkg/dy"
	"galactus/common/middleware/routers"

	"github.com/gin-gonic/gin"
)

var router *routers.GinRouter

func init() {
	router = routers.NewGinRouter()
	InitAllRouters(router)
}

func Run(middleware ...gin.HandlerFunc) error {
	router.Use(middleware...)
	return router.Run()
}

// InitAllRouters 初始化所有router

func InitAllRouters(router *routers.GinRouter) {
	router.Include(device.WebDeviceRouters)
	router.Include(dy.UserRouters)

}
