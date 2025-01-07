package a

import (
	"galactus/blade/internal/initialize"
	"log"
)

func Init() {
	if err := initialize.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
}
