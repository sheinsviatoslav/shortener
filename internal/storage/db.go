package storage

import (
	"database/sql"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"log"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("pgx", *config.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
}
