// Package getuserurls allows to get multiple created urls of current user
package getuserurls

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sheinsviatoslav/shortener/internal/auth"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils"
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
	secretKey, err := hex.DecodeString(common.SecretKey)
	if err != nil {
		http.Error(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ReadEncryptedCookie(r, "userID", secretKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userID == "" {
		userID = utils.GetUserID(r)
	}

	urls, err := h.storage.GetUserUrls(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if _, err := w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
