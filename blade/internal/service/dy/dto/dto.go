package dto

import (
	"galactus/blade/internal/service/device/dto"
	"galactus/blade/internal/service/dy"
)

func NewDyBaseEntity(webDevice *dto.WebDeviceDTO) *dy.DyBaseEntity {
	dyBaseEntity := &dy.DyBaseEntity{
		WebDevice: webDevice,
	}
	return dyBaseEntity
}
