package main

import (
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/sheinsviatoslav/shortener/internal/handlers/geturl"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			createurl.PostHandler(w, req, storage.URLMap)
		} else if req.Method == http.MethodGet {
			geturl.GetHandler(w, req, storage.URLMap)
		} else {
			return
		}
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
