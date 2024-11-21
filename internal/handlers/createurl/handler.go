package createurl

import (
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	DefaultHashLength = 8
)

func Handler(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputURL := string(bodyBytes)
	if inputURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(inputURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items := make([]storage.URLItem, 0)

	if _, err := os.Stat(*config.FileStoragePath); err == nil {
		var fileReader, err = storage.NewConsumer(*config.FileStoragePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer fileReader.Close()

		if fileReader != nil {
			urlItems, err := fileReader.ReadURLItems()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			items = urlItems.Items
		}
	}

	var shortURL string
	var isURLExists bool

	for _, urlItem := range items {
		if urlItem.OriginalURL == inputURL {
			isURLExists = true
			shortURL = urlItem.ShortURL
			break
		}
	}

	if !isURLExists {
		shortURL = hash.Generator(DefaultHashLength)
		var fileWriter, err = storage.NewProducer(*config.FileStoragePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer fileWriter.Close()

		if fileWriter != nil {
			id := 1
			if len(items) > 0 {
				id = items[len(items)-1].ID + 1
			}
			items = append(items, storage.URLItem{
				ID:          id,
				OriginalURL: inputURL,
				ShortURL:    shortURL,
			})

			if err = fileWriter.WriteURLItems(&storage.URLItems{
				Items: items,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	u, _ := url.Parse(*config.BaseURL)
	relative, _ := url.Parse(shortURL)

	w.Header().Set("Content-Type", "text/plain")
	if isURLExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
