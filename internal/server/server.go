package server

import (
	"fmt"
	"net/http"

	"github.com/aswnrj/task-queue-platform/internal/handlers"
	"github.com/aswnrj/task-queue-platform/internal/store"
)

func Run(addr string) error {
	memStore := store.NewMemoryStore()
	taskHandler := handlers.NewTaskHandler(memStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.Get)

	fmt.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, mux)
}
