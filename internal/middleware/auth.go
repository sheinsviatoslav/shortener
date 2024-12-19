package middleware

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/sheinsviatoslav/shortener/internal/auth"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"log"
	"net/http"
)

var CurrentUserID string

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey, err := hex.DecodeString(common.SecretKey)
		if err != nil {
			log.Fatal(err)
			return
		}

		generatedUserID := uuid.New().String()

		cookie := http.Cookie{
			Name:  "userID",
			Value: generatedUserID,
		}

		value, err := auth.ReadEncryptedCookie(r, "userID", secretKey)
		if err != nil {
			if err := auth.WriteEncryptedCookie(w, cookie, secretKey); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			CurrentUserID = generatedUserID
		}

		if value != "" {
			CurrentUserID = value
		}

		next.ServeHTTP(w, r)
	})
}
