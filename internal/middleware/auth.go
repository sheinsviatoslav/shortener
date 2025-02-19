package middleware

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"

	"github.com/sheinsviatoslav/shortener/internal/auth"
	"github.com/sheinsviatoslav/shortener/internal/common"
)

type contextKey string

// UserIDKey is a key for stored userID in context
const UserIDKey contextKey = "userID"

// WithAuth is middleware for checking if user has cookie and write cookie if he doesn't
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey, err := hex.DecodeString(common.SecretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		generatedUserID := uuid.New().String()

		cookie := http.Cookie{
			Name:  "userID",
			Value: generatedUserID,
		}

		var ctx context.Context

		value, err := auth.ReadEncryptedCookie(r, "userID", secretKey)
		if err != nil {
			if err := auth.WriteEncryptedCookie(w, cookie, secretKey); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(r.Context(), UserIDKey, generatedUserID)
		}

		if value != "" {
			ctx = context.WithValue(r.Context(), UserIDKey, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
