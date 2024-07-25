package main

import (
	"Messaggio/internal/api"
	"Messaggio/internal/config"
	"Messaggio/internal/storage"
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	db, err := storage.Connection(cfg)
	if err != nil {
		log.Error("Error DB connection: ", err)
	}

	log.Info("DB connected")
	s := storage.NewStorage(db, log)
	if err = s.CreateTable(); err != nil {
		log.Error("Table create error: ", err)
	}

	router := chi.NewRouter()

	handler := api.NewHandler(s, log)
	handler.Register(router)

	log.Info("Starting server", slog.String("port", cfg.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router); err != nil {
		log.Error("Error starting server: ", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
