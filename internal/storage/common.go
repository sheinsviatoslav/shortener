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

type Storage interface {
	GetOriginalURLByShortURL(string) (string, error)
	GetShortURLByOriginalURL(string) (string, bool, error)
	AddNewURL(string, string) error
	AddManyUrls(InputManyUrls) (OutputManyUrls, error)
}
