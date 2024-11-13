package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"net/url"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

var (
	ServerAddr = flag.String("a", ":8080", "server address")
	BaseURL    = flag.String("b", "http://localhost:8080/", "base address of shortened URL")
)

func Init() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if cfg.ServerAddr != "" {
		*ServerAddr = cfg.ServerAddr
	}

	if cfg.BaseURL != "" {
		*BaseURL = cfg.BaseURL
	}

	if _, err := url.ParseRequestURI(*BaseURL); err != nil {
		log.Fatal(err)
	}
}
