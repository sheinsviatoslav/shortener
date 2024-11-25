package createurl

import (
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"io"
	"net/http"
	"net/url"
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

	urlMap := storage.URLMap{}

	urlItems, err := storage.ReadURLItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if urlItems != nil {
		urlMap = *urlItems
	}

	shortURL, isOriginalURLExists := urlMap[inputURL]

	if !isOriginalURLExists {
		shortURL = hash.Generator(DefaultHashLength)
		urlMap[inputURL] = shortURL
		err := storage.WriteURLItemToFile(urlMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	u, _ := url.Parse(*config.BaseURL)
	relative, _ := url.Parse(shortURL)

	w.Header().Set("Content-Type", "text/plain")
	if isOriginalURLExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
