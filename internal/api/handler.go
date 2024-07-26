package api

import (
	"Messaggio/internal/models"
	"Messaggio/internal/storage"
	"encoding/json"
	"github.com/go-chi/chi"
	"html/template"
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
	router.Get("/", h.formHandler)
	router.Post("/messages", h.messageHandler)
	router.Get("/status", h.statusHandler)
}

func (h *Handler) formHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("api.handler.formHandler")

	tmp, err := template.ParseFiles("/Users/fhyman/GolandProjects/Messaggio/ui/index.html")
	if err != nil {
		http.Error(w, "Error parse file", http.StatusInternalServerError)
		h.logger.Error("error: ", err)
		return
	}

	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing file", http.StatusInternalServerError)
		h.logger.Error("error: ", err)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		h.logger.Error("error: ", err)
		return
	}
	http.Redirect(w, r, "/messages", http.StatusSeeOther)
}

func (h *Handler) messageHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("api.handler.messageHandler")

	//msgContent := r.FormValue("content")
	//if msgContent == "" {
	//	http.Error(w, "Message content is empty", http.StatusBadRequest)
	//	return
	//}

	var msg models.Message
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	} else {
		msgContent := r.FormValue("message")
		if msgContent == "" {
			http.Error(w, "Message content is empty", http.StatusBadRequest)
			return
		}
		msg = models.Message{Content: msgContent}
	}

	err := h.storage.SaveMessage(msg)
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

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshaledStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
