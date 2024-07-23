package api

import (
	"Messaggio/internal/models"
	"Messaggio/internal/storage"
	"encoding/json"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
)

type Handler struct {
	db     *storage.Storage
	logger *slog.Logger
}

func NewHandler(db *storage.Storage, log *slog.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: log,
	}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Post("/messages", h.messageHandler)
	router.Get("/status", h.statusHandler)
}

func (h *Handler) messageHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("api.handler.messageHandler")

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err = h.db.SaveMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) statusHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	marshaledStats, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshaledStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
