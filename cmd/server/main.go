package main

import (
	"flag"
	"log"
	"net/http"

	"sovana/internal/api"
	"sovana/internal/middleware"
	"sovana/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	staticDir := flag.String("static", "web", "directory for static assets")
	flag.Parse()

	store := storage.NewMemoryStorage()
	handler := api.NewHandler(store, *staticDir)

	mux := handler.Routes()
	server := &http.Server{
		Addr:    *addr,
		Handler: middleware.Logger(mux),
	}

	log.Printf("listening on %s", *addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
