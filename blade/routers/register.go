package routers

import (
	"galactus/blade/pkg/device"
	"galactus/blade/pkg/dy"
	"galactus/blade/pkg/ip"
	"galactus/blade/pkg/test"
	"galactus/common/middleware/routers"
)

func registerHandler() []routers.Handler {
	handlers := []routers.Handler{
		dy.NewUserHandler(),
		dy.NewVideoHandler(),
		device.NewWebDeviceHandler(),
		ip.NewIpHandler(),
		test.NewTestHandler(),
	}
	return handlers
}
