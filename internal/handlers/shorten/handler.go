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
)

type ReqBody struct {
	URL string
}

type RespBody struct {
	Result string `json:"result"`
}

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	var reqBody ReqBody
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL := reqBody.URL

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
		shortURL = hash.Generator(createurl.DefaultHashLength)
		if err := h.storage.AddNewURL(originalURL, shortURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
	if isExists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if _, err := w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
