package internalstats

import (
	"encoding/json"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"net"
	"net/http"
)

// Handler is a handler type
type Handler struct {
	storage storage.Storage
}

// NewHandler is a handler constructor
func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// Handle is a main handler method
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if *config.TrustedSubnet == "" {
		http.Error(w, "trusted subnet variable is empty", http.StatusForbidden)
		return
	}

	_, subnet, _ := net.ParseCIDR(*config.TrustedSubnet)
	ip := net.ParseIP(r.Header.Get("X-Real-IP"))

	if subnet.Contains(ip) {
		http.Error(w, "invalid X-Real-IP header", http.StatusForbidden)
		return
	}

	stats, err := h.storage.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
