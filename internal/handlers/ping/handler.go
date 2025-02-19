// Package ping checks if database is successfully connected
package ping

import (
	"database/sql"
	"net/http"
)

// Handler is a handler type
type Handler struct {
	db *sql.DB
}

// NewHandler is a handler constructor
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// Handle is a main handler method
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := h.db.PingContext(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("database is connected")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
