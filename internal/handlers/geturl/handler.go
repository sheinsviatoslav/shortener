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

	var resultURL string

	for _, urlItem := range items {
		if urlItem.ShortURL == id {
			resultURL = urlItem.OriginalURL
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
