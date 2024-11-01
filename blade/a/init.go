package a

import (
	"galactus/common/middleware/db"
	"galactus/common/middleware/vipper"
)

func init() {
	vipper.Init()
	db.Init()
}

func Init() {
}
