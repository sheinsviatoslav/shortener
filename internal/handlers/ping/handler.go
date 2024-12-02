package ping

import (
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	if err := storage.DB.PingContext(req.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("database is connected")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
