package initialization

import (
	"galactus/blade/internal/initialize"
	"galactus/blade/routers"
	"log"
)

func Init() {
	if err := initialize.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	// routers init
	routers.Init()
}
