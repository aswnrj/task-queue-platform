package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aswnrj/task-queue-platform/internal/handlers"
	"github.com/aswnrj/task-queue-platform/internal/store"
	"github.com/aswnrj/task-queue-platform/internal/worker"
)

func Run(addr string) error {
	memStore := store.NewMemoryStore()
	pool := worker.NewPool(3, 100, memStore)
	taskHandler := handlers.NewTaskHandler(memStore, pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.Get)

	srv := &http.Server{Addr: addr, Handler: mux}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("shutdown signal received")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("Server starting on %s", addr)
	return srv.ListenAndServe()
}
