package geturl

import (
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "empty path", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(*config.FileStoragePath); err != nil {
		http.Error(w, "file storage not found", http.StatusNotFound)
		return
	}

	fileReader, err := storage.NewConsumer(*config.FileStoragePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer fileReader.Close()

	urlItems, err := storage.ReadURLItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if urlItems == nil {
		http.Error(w, "empty url items", http.StatusBadRequest)
		return
	}

	var resultURL string

	for originalURL, shortURL := range *urlItems {
		if shortURL == id {
			resultURL = originalURL
			break
		}
	}

	if resultURL == "" {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Location", resultURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	if _, err := w.Write([]byte(resultURL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
