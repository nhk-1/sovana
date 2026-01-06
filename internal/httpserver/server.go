package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"sovana/internal/metrics"
	"sovana/internal/realtime"
	"sovana/internal/storage"
	"sovana/internal/ws"
)

// Server bundles HTTP handlers and dependencies.
type Server struct {
	store  *storage.MemoryStore
	hub    *realtime.Hub
	logger *log.Logger
	server *http.Server
}

// Config drives HTTP server creation.
type Config struct {
	Addr   string
	Logger *log.Logger
}

// New builds a new Server instance.
func New(cfg Config, store *storage.MemoryStore, hub *realtime.Hub) *Server {
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}

	s := &Server{
		store:  store,
		hub:    hub,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", s.loggingMiddleware(http.HandlerFunc(s.handleMetrics)))
	mux.Handle("/ws", s.loggingMiddleware(http.HandlerFunc(s.handleWS)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	s.server = &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return s
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Addr exposes the configured listen address.
func (s *Server) Addr() string { return s.server.Addr }

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var m metrics.Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if m.Timestamp.IsZero() {
		m.Timestamp = time.Now().UTC()
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.store.Add(m)
	s.hub.Broadcast(m)

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrade(w, r)
	if err != nil {
		s.logger.Printf("upgrade websocket: %v", err)
		return
	}

	client := realtime.NewClient(s.hub, conn)
	s.hub.Register(client)

	ctx, cancel := context.WithCancel(r.Context())
	go func() {
		defer cancel()
		// Drain frames to detect client disconnects.
		_ = conn.Drain()
	}()

	client.WritePump(ctx)
	s.hub.Unregister(client)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
