package createurl

import (
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultHashLength = 8
)

func Handler(w http.ResponseWriter, req *http.Request, storage map[string]string) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := string(bodyBytes)
	if body == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	if foundURL, ok := storage[body]; ok {
		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(foundURL)

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		hashVal := hash.Generator(DefaultHashLength)
		storage[body] = hashVal
		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(hashVal)

		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write([]byte(u.ResolveReference(relative).String())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
