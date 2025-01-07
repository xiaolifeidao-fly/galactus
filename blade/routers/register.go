package routers

import (
	"galactus/blade/pkg/device"
	"galactus/blade/pkg/dy"
	"galactus/common/middleware/routers"
)

func registerHandler() []routers.Handler {
	handlers := []routers.Handler{
		dy.NewUserHandler(),
		dy.NewVideoHandler(),
		device.NewWebDeviceHandler(),
	}
	return handlers
}
