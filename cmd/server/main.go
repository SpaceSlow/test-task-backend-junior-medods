package main

import (
	"log"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal"
)

func main() {
	if err := internal.RunServer(); err != nil {
		log.Printf("error: %s", err)
	}
}
