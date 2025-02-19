// Package shortenbatch allows to create multiple short urls from multiple original urls
package shortenbatch

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils"
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
	var reqBody storage.InputManyUrls
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(reqBody) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	respBody, err := h.storage.AddManyUrls(r.Context(), reqBody, utils.GetUserID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
