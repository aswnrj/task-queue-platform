package main

import (
	"log"

	"github.com/aswnrj/task-queue-platform/internal/server"
)

func main() {
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
