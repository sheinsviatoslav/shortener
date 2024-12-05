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

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
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

	shortURL, isExists, storageErr := h.storage.GetShortURLByOriginalURL(originalURL)
	if storageErr != nil {
		http.Error(w, storageErr.Error(), http.StatusInternalServerError)
		return
	}

	if !isExists {
		shortURL = hash.Generator(DefaultHashLength)
		if err := h.storage.AddNewURL(originalURL, shortURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	u, _ := url.Parse(*config.BaseURL)
	relative, _ := url.Parse(shortURL)

	w.Header().Set("Content-Type", "text/plain")
	if isExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
