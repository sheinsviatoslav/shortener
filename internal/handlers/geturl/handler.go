package geturl

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sheinsviatoslav/shortener/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	if shortURL == "" {
		http.Error(w, "empty path", http.StatusBadRequest)
		return
	}

	originalURL, isDeleted, err := h.storage.GetOriginalURLByShortURL(r.Context(), shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isDeleted {
		http.Error(w, "url is already deleted", http.StatusGone)
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
