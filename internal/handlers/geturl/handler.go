package geturl

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request, storage map[string]string) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "empty path", http.StatusBadRequest)
		return
	}

	var resultURL string
	for k, v := range storage {
		if v == id {
			resultURL = k
			break
		}
	}

	if resultURL == "" {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Add("Location", resultURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	if _, err := w.Write([]byte(resultURL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
