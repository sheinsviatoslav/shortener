package shorten

import (
	"bytes"
	"encoding/json"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"net/http"
	"net/url"
	"os"
)

type ReqBody struct {
	URL string
}

type RespBody struct {
	Result string `json:"result"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var reqBody ReqBody
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(reqBody.URL); err != nil {
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
	isURLExists := false

	for _, urlItem := range items {
		if urlItem.OriginalURL == reqBody.URL {
			isURLExists = true
			shortURL = urlItem.ShortURL
			break
		}
	}

	if !isURLExists {
		shortURL = hash.Generator(createurl.DefaultHashLength)
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
				OriginalURL: reqBody.URL,
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

	respBody := RespBody{
		Result: u.ResolveReference(relative).String(),
	}

	jsonResp, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if isURLExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
