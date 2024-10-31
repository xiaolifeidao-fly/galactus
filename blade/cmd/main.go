package main

import (
	"galactus/common/middleware/db"
	"galactus/common/middleware/vipper"
)

func Init() {
	vipper.Init()
	db.Init()
}
