package main

import (
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/routes"
	"net/http"
)

func main() {
	config.Init()
	err := http.ListenAndServe(*config.ServerAddr, routes.MainRouter())
	if err != nil {
		panic(err)
	}
}
