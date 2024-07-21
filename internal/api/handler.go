package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Register(router *chi.Mux) {
	router.Post("/messages", messageHandler)
	router.Get("/status", statusHandler)
}

func messageHandler(w http.ResponseWriter, r *http.Request) {

}

func statusHandler(w http.ResponseWriter, r *http.Request) {

}
