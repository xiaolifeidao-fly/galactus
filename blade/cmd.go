package main

import (
	"galactus/blade/a"
	"galactus/blade/routers"
)

func main() {
	a.Init()
	routers.Run()
}
