package storage

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
	GetOriginalURLByShortURL(string) (string, bool, error)
	GetShortURLByOriginalURL(string) (string, bool, error)
	AddNewURL(string, string, string) error
	AddManyUrls(InputManyUrls, string) (OutputManyUrls, error)
	GetUserUrls(string) (UserUrls, error)
	DeleteUserUrls([]string, string) error
}
