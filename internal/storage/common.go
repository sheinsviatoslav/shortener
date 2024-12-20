package storage

import "context"

type InputManyUrlsItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type InputManyUrls []InputManyUrlsItem

type OutputManyUrlsItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
type OutputManyUrls []OutputManyUrlsItem

type UserUrlsItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserUrls []UserUrlsItem

type Storage interface {
	GetOriginalURLByShortURL(context.Context, string) (string, bool, error)
	GetShortURLByOriginalURL(context.Context, string) (string, bool, error)
	AddNewURL(context.Context, string, string, string) error
	AddManyUrls(context.Context, InputManyUrls, string) (OutputManyUrls, error)
	GetUserUrls(context.Context, string) (UserUrls, error)
	DeleteUserUrls(context.Context, []string, string) error
}
