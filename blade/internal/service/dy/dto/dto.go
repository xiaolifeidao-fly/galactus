package dto

import (
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dy"
)

func NewDyBaseEntity(webDevice *dto.WebDeviceDTO, ip string) *dy.DyBaseEntity {
	dyBaseEntity := &dy.DyBaseEntity{
		WebDevice: webDevice,
		Ip:        ip,
	}
	return dyBaseEntity
}
