package main

import (
	"galactus/blade/initialization"
	"galactus/blade/routers"
)

func main() {
	initialization.Init()
	routers.Run()
}
