package routes

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/ping"
	"github.com/sheinsviatoslav/shortener/internal/handlers/shorten"
	"github.com/sheinsviatoslav/shortener/internal/middleware"
	"time"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogger)
	r.Use(middleware.GzipHandle)
	r.Use(chiMiddleware.Timeout(1000 * time.Millisecond))

	r.Get("/{id}", geturl.Handler)
	r.Post("/", createurl.Handler)
	r.Post("/api/shorten", shorten.Handler)
	r.Get("/ping", ping.Handler)

	return r
}
