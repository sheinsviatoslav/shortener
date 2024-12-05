package storage

import (
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

func (m *MemStorage) AddNewURL(originalURL string, shortURL string) error {
	m.m.Lock()
	defer m.m.Unlock()
	m.data[originalURL] = shortURL

	return nil
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
