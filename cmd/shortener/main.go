package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ReqBody struct {
	Url string
}

func mainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		var reqBody ReqBody
		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if reqBody.Url == "" {
			http.Error(res, "url is required", http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(reqBody.Url))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if req.URL.Path == "/" {
			http.Error(res, "Invalid path", http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Add("Location", "https://practicum.yandex.ru/")
		res.WriteHeader(http.StatusTemporaryRedirect)
		_, err := res.Write([]byte(req.URL.Path))
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
