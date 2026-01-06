package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sovana/internal/api"
	"sovana/internal/middleware"
	"sovana/internal/storage"
)

func main() {
	addr := envOrDefault("ADDR", ":8080")

	store := storage.NewMemoryStore()
	apiServer := &api.Server{Store: store}

	mux := http.NewServeMux()
	apiServer.RegisterRoutes(mux)

	webRoot := resolveWebRoot()
	fileServer := http.FileServer(http.Dir(webRoot))
	mux.Handle("/", fileServer)

	handler := middleware.Logging(mux)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func resolveWebRoot() string {
	if _, err := os.Stat("web"); err == nil {
		return "web"
	}

	if _, err := os.Stat("app/web"); err == nil {
		return "app/web"
	}

	return "web"
}
