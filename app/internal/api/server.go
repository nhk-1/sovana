package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"sovana/internal/model"
	"sovana/internal/storage"
)

// Server holds dependencies for HTTP handlers.
type Server struct {
	Store storage.Store
}

// RegisterRoutes attaches API routes to the provided mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/items", s.itemsHandler)
	mux.HandleFunc("/api/items/", s.itemHandler)
}

func (s *Server) itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleList(w, r)
	case http.MethodPost:
		s.handleCreate(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *Server) itemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/items/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := s.Store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	items, err := s.Store.List(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	payload.Title = strings.TrimSpace(payload.Title)
	payload.Description = strings.TrimSpace(payload.Description)

	if payload.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	item := model.Item{
		Title:       payload.Title,
		Description: payload.Description,
	}

	created, err := s.Store.Create(r.Context(), item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, http.StatusCreated, created)
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
