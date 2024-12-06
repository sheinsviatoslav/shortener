package storage

import (
	"database/sql"
	"errors"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"net/url"
)

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
			"id uuid PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
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
	if _, err := p.DB.Exec(`INSERT INTO urls (original_url, short_url) VALUES($1, $2)`, originalURL, shortURL); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) AddManyUrls(urls InputManyUrls) (OutputManyUrls, error) {
	var output OutputManyUrls
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	for _, item := range urls {
		if item.OriginalURL == "" {
			tx.Rollback()
			return nil, errors.New("url is required")
		}

		if _, err = url.ParseRequestURI(item.OriginalURL); err != nil {
			tx.Rollback()
			return nil, err
		}

		shortURL, isExists, dbErr := p.GetShortURLByOriginalURL(item.OriginalURL)
		if dbErr != nil {
			tx.Rollback()
			return nil, dbErr
		}

		if !isExists {
			shortURL = hash.Generator(common.DefaultHashLength)
			if _, err = tx.Exec(`INSERT INTO urls (original_url, short_url) VALUES($1, $2)`, item.OriginalURL, shortURL); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(shortURL)

		output = append(output, OutputManyUrlsItem{CorrelationID: item.CorrelationID, ShortURL: u.ResolveReference(relative).String()})
	}

	tx.Commit()
	return output, nil
}
