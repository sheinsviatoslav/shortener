package shorten

import (
	"bytes"
	"encoding/json"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
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

func Handler(w http.ResponseWriter, r *http.Request, storage map[string]string) {
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

	id, isURLExists := storage[reqBody.URL]

	if !isURLExists {
		hashVal := hash.Generator(createurl.DefaultHashLength)
		storage[reqBody.URL] = hashVal
		id = hashVal
	}

	u, _ := url.Parse(*config.BaseURL)
	relative, _ := url.Parse(id)

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
