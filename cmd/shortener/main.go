package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/routes"
)

func main() {
	config.Init()

	log.Println("listen on", *config.ServerAddr)
	log.Fatal(http.ListenAndServe(*config.ServerAddr, routes.MainRouter()))
}
