package storage

import (
	"database/sql"
	"errors"
	"github.com/sheinsviatoslav/shortener/internal/config"
)

type Storage interface {
	GetOriginalURLByShortURL(string) (string, error)
	GetShortURLByOriginalURL(string) (string, bool, error)
	AddNewURL(string, string) error
}

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage() *PgStorage {
	return &PgStorage{
		DB: nil,
	}
}

func (p *PgStorage) Connect() error {
	var err error
	p.DB, err = sql.Open("pgx", *config.DatabaseDSN)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(
		"CREATE TABLE IF NOT EXISTS urls (" +
			"id SERIAL PRIMARY KEY, " +
			"original_url TEXT NOT NULL UNIQUE, " +
			"short_url TEXT NOT NULL UNIQUE)")
	if err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetOriginalURLByShortURL(shortURL string) (string, error) {
	row := p.DB.QueryRow(`SELECT original_url FROM urls WHERE short_url = $1`, shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}

	return originalURL, nil
}

func (p *PgStorage) GetShortURLByOriginalURL(originalURL string) (string, bool, error) {
	row := p.DB.QueryRow(`SELECT short_url FROM urls WHERE original_url = $1`, originalURL)
	var shortURL string

	if err := row.Scan(&shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}

		return "", false, err
	}

	return shortURL, true, nil
}

func (p *PgStorage) AddNewURL(originalURL string, shortURL string) error {
	_, err := p.DB.Exec(`INSERT INTO urls (original_url, short_url) VALUES($1, $2)`, originalURL, shortURL)
	if err != nil {
		return err
	}

	return nil
}
