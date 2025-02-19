package routes

import (
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/deleteuserurls"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/getuserurls"
	"github.com/sheinsviatoslav/shortener/internal/handlers/ping"
	"github.com/sheinsviatoslav/shortener/internal/handlers/shorten"
	"github.com/sheinsviatoslav/shortener/internal/handlers/shortenbatch"
	"github.com/sheinsviatoslav/shortener/internal/middleware"
	"github.com/sheinsviatoslav/shortener/internal/storage"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithLogger)
	r.Use(middleware.GzipHandle)
	r.Use(middleware.WithAuth)
	r.Use(chiMiddleware.Timeout(1000 * time.Millisecond))

	var st storage.Storage
	if *config.DatabaseDSN != "" {
		pgStorage := storage.NewPgStorage()
		if err := pgStorage.Connect(); err != nil {
			log.Fatal(err)
		}

		st = pgStorage
		r.Get("/ping", ping.NewHandler(pgStorage.DB).Handle)
	} else if *config.FileStoragePath != "" {
		st = storage.NewFileStorage()
	} else {
		st = storage.NewMemStorage()
	}

	r.Get("/{shortURL}", geturl.NewHandler(st).Handle)
	r.Post("/", createurl.NewHandler(st).Handle)
	r.Post("/api/shorten", shorten.NewHandler(st).Handle)
	r.Post("/api/shorten/batch", shortenbatch.NewHandler(st).Handle)
	r.Get("/api/user/urls", getuserurls.NewHandler(st).Handle)
	r.Delete("/api/user/urls", deleteuserurls.NewHandler(st).Handle)
	r.Mount("/debug", chiMiddleware.Profiler())

	return r
}
