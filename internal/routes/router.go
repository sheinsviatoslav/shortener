package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/shorten"
	"github.com/sheinsviatoslav/shortener/internal/middleware"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogger)
	r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
		geturl.Handler(w, req, storage.URLMap)
	})
	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		createurl.Handler(w, req, storage.URLMap)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, req *http.Request) {
		shorten.Handler(w, req, storage.URLMap)
	})

	return r
}
