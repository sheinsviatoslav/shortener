package middleware

import (
	"context"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/sheinsviatoslav/shortener/internal/auth"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

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
