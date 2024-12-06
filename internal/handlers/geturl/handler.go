package geturl

import (
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
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
	shortURL := chi.URLParam(req, "shortURL")
	if shortURL == "" {
		http.Error(w, "empty path", http.StatusBadRequest)
		return
	}

	originalURL, err := h.storage.GetOriginalURLByShortURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if originalURL == "" {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	if _, err := w.Write([]byte(originalURL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
