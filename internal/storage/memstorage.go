package storage

import (
	"errors"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"net/url"
	"sync"
)

type MemStorage struct {
	data map[string]string
	m    sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func (m *MemStorage) GetShortURLByOriginalURL(originalURL string) (string, bool, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if shortURL, ok := m.data[originalURL]; ok {
		return shortURL, true, nil
	}

	return "", false, nil
}

func (m *MemStorage) GetOriginalURLByShortURL(inputShortURL string) (string, error) {
	m.m.Lock()
	defer m.m.Unlock()

	for originalURL, shortURL := range m.data {
		if shortURL == inputShortURL {
			return originalURL, nil
		}
	}

	return "", nil
}

func (m *MemStorage) AddNewURL(originalURL string, shortURL string, _ string) error {
	m.m.Lock()
	defer m.m.Unlock()
	m.data[originalURL] = shortURL

	return nil
}

func (m *MemStorage) AddManyUrls(urls InputManyUrls, _ string) (OutputManyUrls, error) {
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

func (m *MemStorage) GetUserUrls(_ string) (UserUrls, error) {
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
