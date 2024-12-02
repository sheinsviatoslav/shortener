package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/routes"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	config.Init()
	storage.ConnectDB()

	defer storage.DB.Close()

	log.Println("listen on", *config.ServerAddr)
	log.Fatal(http.ListenAndServe(*config.ServerAddr, routes.MainRouter()))
}
