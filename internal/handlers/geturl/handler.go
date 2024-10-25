package geturl

import (
	"net/http"
)

func GetHandler(w http.ResponseWriter, req *http.Request, storage map[string]string) {
	if req.URL.Path == "/" {
		http.Error(w, "empty path", http.StatusBadRequest)
		return
	}

	var resultURL string
	for k, v := range storage {
		if v == req.URL.Path {
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
	_, err := w.Write([]byte(resultURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
