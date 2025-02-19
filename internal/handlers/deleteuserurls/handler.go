// Package deleteuserurls allows to delete multiple urls
package deleteuserurls

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/sheinsviatoslav/shortener/internal/auth"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/storage"
)

type ReqBody []string

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	secretKey, err := hex.DecodeString(common.SecretKey)
	if err != nil {
		http.Error(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ReadEncryptedCookie(r, "userID", secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	go func(urls ReqBody, id string) {
		_ = h.storage.DeleteUserUrls(r.Context(), reqBody, id)
	}(reqBody, userID)

	if err = h.storage.DeleteUserUrls(r.Context(), reqBody, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
