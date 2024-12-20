package utils

import (
	"github.com/sheinsviatoslav/shortener/internal/middleware"
	"net/http"
)

func GetUserID(r *http.Request) string {
	userID := r.Context().Value(middleware.UserIDKey)

	if userID == nil {
		return ""
	}

	return userID.(string)
}
