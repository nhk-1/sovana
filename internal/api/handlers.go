package api

import (
	"encoding/json"
	"log"
	"net/http"

	"sovana/internal/detector"
	"sovana/internal/parser"
	"sovana/internal/storage"
)

// Handler wires storage and business logic for HTTP endpoints.
type Handler struct {
	Store       *storage.MemoryStorage
	StaticDir   string
	MaxUploadMB int64
}

func NewHandler(store *storage.MemoryStorage, staticDir string) *Handler {
	return &Handler{Store: store, StaticDir: staticDir, MaxUploadMB: 5}
}

// Routes registers all HTTP handlers.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", h.health)
	mux.HandleFunc("/api/upload", h.upload)
	mux.HandleFunc("/api/subscriptions", h.subscriptions)

	fileServer := http.FileServer(http.Dir(h.StaticDir))
	mux.Handle("/", fileServer)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(h.MaxUploadMB << 20); err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("received upload: %s (%d bytes)", header.Filename, header.Size)
	transactions, err := parser.ParseCSV(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subs := detector.DetectSubscriptions(transactions)
	h.Store.Save(transactions, subs)
	respondJSON(w, http.StatusOK, map[string]interface{}{"subscriptions": subs, "count": len(subs)})
}

func (h *Handler) subscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	respondJSON(w, http.StatusOK, h.Store.Subscriptions())
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
