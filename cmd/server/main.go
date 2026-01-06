package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sovana/internal/httpserver"
	"sovana/internal/realtime"
	"sovana/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "server ", log.LstdFlags)

	store := storage.NewMemoryStore()
	hub := realtime.NewHub()

	srv := httpserver.New(httpserver.Config{Addr: ":8080", Logger: logger}, store, hub)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go hub.Run(ctx)

	go func() {
		logger.Printf("listening on %s", srvAddr(srv))
		if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
	}
}

func srvAddr(s *httpserver.Server) string {
	return s.Addr()
}
