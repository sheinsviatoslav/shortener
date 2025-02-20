// Package createurl allows to create short url from original url using plain text content type
package createurl

import (
	"io"
	"net/http"
	"net/url"

	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
)

// Handler is a handler type
type Handler struct {
	storage storage.Storage
}

// NewHandler is a handler constructor
func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// Handle is a main handler method
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL := string(bodyBytes)
	if originalURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, isExists, storageErr := h.storage.GetShortURLByOriginalURL(r.Context(), originalURL)
	if storageErr != nil {
		http.Error(w, storageErr.Error(), http.StatusInternalServerError)
		return
	}

	if !isExists {
		shortURL = hash.Generator(common.DefaultHashLength)
		if err := h.storage.AddNewURL(r.Context(), originalURL, shortURL, utils.GetUserID(r)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	u, _ := url.Parse(*config.BaseURL)
	relative, _ := url.Parse(shortURL)

	w.Header().Set("Content-Type", "text/plain")
	if isExists {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
