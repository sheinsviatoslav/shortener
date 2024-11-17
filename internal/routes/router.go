package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/middleware"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", middleware.WithLogger(func(w http.ResponseWriter, req *http.Request) {
		geturl.Handler(w, req, storage.URLMap)
	}))
	r.Post("/", middleware.WithLogger(func(w http.ResponseWriter, req *http.Request) {
		createurl.Handler(w, req, storage.URLMap)
	}))

	return r
}
