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
	storage *storage.Storage
	//producer *kafka.Producer
	logger *slog.Logger
}

// func NewHandler(db *storage.Storage, log *slog.Logger, prod *kafka.Producer) *Handler {
func NewHandler(db *storage.Storage, log *slog.Logger) *Handler {
	return &Handler{
		storage: db,
		//producer: prod,
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
		return
	}

	err = h.storage.SaveMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) statusHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.storage.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	marshaledStats, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshaledStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
