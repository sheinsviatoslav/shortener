package storage

import (
	"context"
	"errors"
	"net/url"
	"slices"
	"sync"

	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
)

// MemStorage is a storage type
type MemStorage struct {
	data map[string]string
	m    sync.Mutex
}

// NewMemStorage constructs new MemStorage struct
func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

// GetOriginalURLByShortURL returns original url
func (m *MemStorage) GetOriginalURLByShortURL(_ context.Context, inputShortURL string) (string, bool, error) {
	m.m.Lock()
	defer m.m.Unlock()

	for originalURL, shortURL := range m.data {
		if shortURL == inputShortURL {
			return originalURL, false, nil
		}
	}

	return "", false, nil
}

// GetShortURLByOriginalURL returns short url
func (m *MemStorage) GetShortURLByOriginalURL(_ context.Context, originalURL string) (string, bool, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if shortURL, ok := m.data[originalURL]; ok {
		return shortURL, true, nil
	}

	return "", false, nil
}

// AddNewURL adds originalURL-shortURL pair into the storage
func (m *MemStorage) AddNewURL(_ context.Context, originalURL string, shortURL string, _ string) error {
	m.m.Lock()
	defer m.m.Unlock()
	m.data[originalURL] = shortURL

	return nil
}

// AddManyUrls adds multiple originalURL-shortURL pairs into the storage
func (m *MemStorage) AddManyUrls(_ context.Context, urls InputManyUrls, _ string) (OutputManyUrls, error) {
	m.m.Lock()
	defer m.m.Unlock()
	var output OutputManyUrls

	for _, item := range urls {
		if item.OriginalURL == "" {
			return nil, errors.New("url is required")
		}

		if _, err := url.ParseRequestURI(item.OriginalURL); err != nil {
			return nil, err
		}

		shortURL, isExists := m.data[item.OriginalURL]

		if !isExists {
			shortURL = hash.Generator(common.DefaultHashLength)
			m.data[item.OriginalURL] = shortURL
		}

		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(shortURL)

		output = append(output, OutputManyUrlsItem{CorrelationID: item.CorrelationID, ShortURL: u.ResolveReference(relative).String()})
	}

	return output, nil
}

// GetUserUrls returns multiple originalURL-shortURL pairs of current user
func (m *MemStorage) GetUserUrls(_ context.Context, _ string) (UserUrls, error) {
	m.m.Lock()
	defer m.m.Unlock()

	output := make(UserUrls, 0)
	for originalURL, shortURL := range m.data {
		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(shortURL)
		output = append(output, UserUrlsItem{OriginalURL: originalURL, ShortURL: u.ResolveReference(relative).String()})
	}

	return output, nil
}

// DeleteUserUrls deletes multiple short urls
func (m *MemStorage) DeleteUserUrls(_ context.Context, shortUrls []string, _ string) error {
	m.m.Lock()
	defer m.m.Unlock()

	for originalURL, shortURL := range m.data {
		if slices.Contains(shortUrls, shortURL) {
			delete(m.data, originalURL)
		}
	}

	return nil
}
