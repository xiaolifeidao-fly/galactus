package repository

import (
	"galactus/common/middleware/db"
)

type WebDevice struct {
	db.BaseEntity
}

func (w *WebDevice) TableName() string {
	return "web_device"
}

type WebDeviceRepository struct {
	db.Repository[*WebDevice]
}
