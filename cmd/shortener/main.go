package main

import (
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/routes"
	"log"
	"net/http"
)

func main() {
	config.Init()

	log.Println("listen on", *config.ServerAddr)
	log.Fatal(http.ListenAndServe(*config.ServerAddr, routes.MainRouter()))
}
