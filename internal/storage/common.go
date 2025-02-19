package storage

import "context"

// InputManyUrlsItem is type for input item in AddManyUrls function
type InputManyUrlsItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// InputManyUrls is type for input in AddManyUrls function
type InputManyUrls []InputManyUrlsItem

// OutputManyUrlsItem is type for output item in AddManyUrls function
type OutputManyUrlsItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// OutputManyUrls is type for output in AddManyUrls function
type OutputManyUrls []OutputManyUrlsItem

// UserUrlsItem is type for output item in GetUserUrls function
type UserUrlsItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// UserUrls is type for output in GetUserUrls function
type UserUrls []UserUrlsItem

// Storage is the interface for api functions
type Storage interface {
	GetOriginalURLByShortURL(context.Context, string) (string, bool, error)
	GetShortURLByOriginalURL(context.Context, string) (string, bool, error)
	AddNewURL(context.Context, string, string, string) error
	AddManyUrls(context.Context, InputManyUrls, string) (OutputManyUrls, error)
	GetUserUrls(context.Context, string) (UserUrls, error)
	DeleteUserUrls(context.Context, []string, string) error
}
