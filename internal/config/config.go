package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"net/url"
	"os"
)

// Config is a config params type
type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS" json:"server_addr"`
	BaseURL         string `env:"BASE_URL" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     string `env:"ENABLE_HTTPS" json:"enable_https"`
	ConfigFile      string `env:"CONFIG_FILE" json:"config_file"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// App config params
var (
	ServerAddr      = flag.String("a", ":8080", "server address")
	BaseURL         = flag.String("b", "http://localhost:8080/", "base address of shortened URL")
	FileStoragePath = flag.String("f", "url_storage.json", "file storage path")
	DatabaseDSN     = flag.String("d", "", "database data source name")
	EnableHTTPS     = flag.String("s", "", "activate https connection")
	ConfigFile      = flag.String("c", "", "config file")
	TrustedSubnet   = flag.String("t", "", "trusted subnet")
)

// Init is a function that checks if config params are from flags or from environment variables
func Init() {
	var cfg Config
	var fileCfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if cfg.ConfigFile != "" {
		*ConfigFile = cfg.ConfigFile
	}

	if *ConfigFile != "" {
		fileConfig, err := os.ReadFile(*ConfigFile)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}

		if len(fileConfig) > 0 {
			if err = json.Unmarshal(fileConfig, &fileCfg); err != nil {
				log.Fatal(err)
			}
		}
	}

	if cfg.ServerAddr != "" {
		*ServerAddr = cfg.ServerAddr
	} else if *ServerAddr == "" && fileCfg.ServerAddr != "" {
		*ServerAddr = fileCfg.ServerAddr
	}

	if cfg.BaseURL != "" {
		*BaseURL = cfg.BaseURL
	} else if *BaseURL == "" && fileCfg.BaseURL != "" {
		*BaseURL = fileCfg.BaseURL
	}

	if cfg.FileStoragePath != "" {
		*FileStoragePath = cfg.FileStoragePath
	} else if *FileStoragePath == "" && fileCfg.FileStoragePath != "" {
		*FileStoragePath = fileCfg.FileStoragePath
	}

	if cfg.DatabaseDSN != "" {
		*DatabaseDSN = cfg.DatabaseDSN
	} else if *DatabaseDSN == "" && fileCfg.DatabaseDSN != "" {
		*DatabaseDSN = fileCfg.DatabaseDSN
	}

	if cfg.EnableHTTPS != "" {
		*EnableHTTPS = cfg.EnableHTTPS
	} else if *EnableHTTPS == "" && fileCfg.EnableHTTPS != "" {
		*EnableHTTPS = fileCfg.EnableHTTPS
	}

	if cfg.TrustedSubnet != "" {
		*TrustedSubnet = cfg.TrustedSubnet
	} else if *TrustedSubnet == "" && fileCfg.TrustedSubnet != "" {
		*TrustedSubnet = fileCfg.TrustedSubnet
	}

	if _, err := url.ParseRequestURI(*BaseURL); err != nil {
		log.Fatal(err)
	}
}
