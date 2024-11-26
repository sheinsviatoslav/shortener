package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/shorten"
	"github.com/sheinsviatoslav/shortener/internal/middleware"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogger)
	r.Use(middleware.GzipHandle)

	r.Get("/{id}", geturl.Handler)
	r.Post("/", createurl.Handler)
	r.Post("/api/shorten", shorten.Handler)

	return r
}
