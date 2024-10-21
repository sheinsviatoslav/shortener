package main

import (
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	"io"
	"net/http"
	"net/url"
)

var urlMap = map[string]string{}
var defaultHashLength = 8

func mainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		body := string(bodyBytes)

		if body == "" {
			http.Error(res, "url is required", http.StatusBadRequest)
			return
		}

		u, err := url.ParseRequestURI(body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		u.Scheme = "http"
		u.Host = req.Host
		res.Header().Set("Content-Type", "text/plain")

		if foundURL, ok := urlMap[body]; ok {
			u.Path = foundURL
			res.WriteHeader(http.StatusOK)
			_, err = res.Write([]byte(u.String()))
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			u.Path = req.URL.Path + hash.Generator(defaultHashLength)
			urlMap[body] = u.Path

			res.WriteHeader(http.StatusCreated)
			_, err = res.Write([]byte(u.String()))
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		if req.URL.Path == "/" {
			http.Error(res, "Invalid path", http.StatusBadRequest)
			return
		}

		var resultURL string
		for k, v := range urlMap {
			if v == req.URL.Path {
				resultURL = k
				break
			}
		}

		if resultURL == "" {
			http.Error(res, "invalid URL path", http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Add("Location", resultURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
		_, err := res.Write([]byte(resultURL))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
