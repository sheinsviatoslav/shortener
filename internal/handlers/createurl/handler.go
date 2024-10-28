package createurl

import (
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"io"
	"net/http"
	"net/url"
)

var defaultHashLength = 8

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

	u, err := url.ParseRequestURI(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Scheme = "http"
	u.Host = req.Host
	w.Header().Set("Content-Type", "text/plain")

	if foundURL, ok := storage[body]; ok {
		u.Path = foundURL
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(u.String()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		hashVal := hash.Generator(defaultHashLength)
		u.Path = hashVal
		storage[body] = hashVal

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(u.String()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
