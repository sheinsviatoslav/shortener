package utils

import (
	"net/http"

	"github.com/sheinsviatoslav/shortener/internal/middleware"
)

func GetUserID(r *http.Request) string {
	userID := r.Context().Value(middleware.UserIDKey)

	if userID == nil {
		return ""
	}

	return userID.(string)
}
