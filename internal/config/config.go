package config

import (
	"flag"
	"log"
	"net/url"

	"github.com/caarlos0/env/v6"
)

// Config is a config params type
type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	EnableHTTPS     string `env:"ENABLE_HTTPS"`
}

// App config params
var (
	ServerAddr      = flag.String("a", ":8080", "server address")
	BaseURL         = flag.String("b", "http://localhost:8080/", "base address of shortened URL")
	FileStoragePath = flag.String("f", "url_storage.json", "file storage path")
	DatabaseDSN     = flag.String("d", "", "database data source name")
	EnableHTTPS     = flag.String("s", "", "activate https connection")
)

// Init is a function that checks if config params are from flags or from environment variables
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

	if cfg.FileStoragePath != "" {
		*FileStoragePath = cfg.FileStoragePath
	}

	if cfg.DatabaseDSN != "" {
		*DatabaseDSN = cfg.DatabaseDSN
	}

	if cfg.EnableHTTPS != "" {
		*EnableHTTPS = cfg.EnableHTTPS
	}

	if _, err := url.ParseRequestURI(*BaseURL); err != nil {
		log.Fatal(err)
	}
}
