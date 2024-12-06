package ping

import (
	"database/sql"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	if err := h.db.PingContext(req.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("database is connected")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
