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
			"id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(), " +
			"original_url TEXT NOT NULL," +
			"short_url TEXT NOT NULL," +
			"user_id UUID NOT NULL," +
			"is_deleted BOOLEAN DEFAULT FALSE)",
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) GetOriginalURLByShortURL(shortURL string) (string, bool, error) {
	row := p.DB.QueryRow(`SELECT original_url, is_deleted FROM urls WHERE short_url = $1`, shortURL)
	var data struct {
		OriginalURL string `json:"original_url"`
		IsDeleted   bool   `json:"is_deleted"`
	}
	if err := row.Scan(&data.OriginalURL, &data.IsDeleted); err != nil {
		return "", false, err
	}

	if data.IsDeleted {
		return "", true, nil
	}

	return data.OriginalURL, false, nil
}

func (p *PgStorage) GetShortURLByOriginalURL(originalURL string) (string, bool, error) {
	query := `SELECT short_url FROM urls WHERE original_url = $1`
	row := p.DB.QueryRow(query, originalURL)
	var shortURL string

	if err := row.Scan(&shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}

		return "", false, err
	}

	return shortURL, true, nil
}

func (p *PgStorage) AddNewURL(originalURL string, shortURL string, userID string) error {
	query := `INSERT INTO urls (original_url, short_url, user_id) VALUES($1, $2, $3)`
	if _, err := p.DB.Exec(query, originalURL, shortURL, userID); err != nil {
		return err
	}

	return nil
}

func (p *PgStorage) AddManyUrls(urls InputManyUrls, userID string) (OutputManyUrls, error) {
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
			query := `INSERT INTO urls (original_url, short_url, user_id) VALUES($1, $2, $3)`
			if _, err = tx.Exec(query, item.OriginalURL, shortURL, userID); err != nil {
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

func (p *PgStorage) GetUserUrls(userID string) (UserUrls, error) {
	query := `SELECT original_url, short_url FROM urls WHERE user_id = $1`
	rows, err := p.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	output := make(UserUrls, 0)

	for rows.Next() {
		var urlItem UserUrlsItem
		err = rows.Scan(&urlItem.OriginalURL, &urlItem.ShortURL)
		if err != nil {
			return nil, err
		}

		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(urlItem.ShortURL)

		output = append(output, UserUrlsItem{OriginalURL: urlItem.OriginalURL, ShortURL: u.ResolveReference(relative).String()})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (p *PgStorage) DeleteUserUrls(shortUrls []string, userID string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	for _, shortURL := range shortUrls {

		query := `UPDATE urls SET is_deleted = true WHERE short_url = $1 AND user_id = $2`
		if _, err = tx.Exec(query, shortURL, userID); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
