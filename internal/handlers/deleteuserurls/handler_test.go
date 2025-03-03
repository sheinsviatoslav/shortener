package deleteuserurls

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/mocks"
)

func TestShortenHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStorage(ctrl)

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name  string
		input []string
		want  want
	}{
		{
			name:  "success",
			input: []string{"s6oGdMVH", "aqZHfy0m"},
			want: want{
				code:        202,
				response:    "",
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.input)
			r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
			w := httptest.NewRecorder()

			m := storage.NewMemStorage()
			assert.NoError(t, m.AddNewURL(r.Context(), "https://yandex.ru/", "3LbIJLJ5", ""))
			assert.NoError(t, m.AddNewURL(r.Context(), "https://practicum.ru/", "99XGYq4c", ""))

			userID := uuid.New().String()
			secretKey, _ := hex.DecodeString(common.SecretKey)
			aesBlock, _ := aes.NewCipher(secretKey)
			aesGCM, _ := cipher.NewGCM(aesBlock)

			nonce := make([]byte, aesGCM.NonceSize())
			io.ReadFull(rand.Reader, nonce)

			plaintext := fmt.Sprintf("%s:%s", "userID", userID)
			encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

			r.AddCookie(&http.Cookie{
				Name:  "userID",
				Value: base64.URLEncoding.EncodeToString(encryptedValue),
			})

			NewHandler(m).Handle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			url, _, _ := m.GetOriginalURLByShortURL(r.Context(), test.input[0])
			assert.Equal(t, "", url)
			url, _, _ = m.GetOriginalURLByShortURL(r.Context(), test.input[1])
			assert.Equal(t, "", url)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	for _, test := range tests {
		t.Run(test.name+" filestorage", func(t *testing.T) {
			body, _ := json.Marshal(test.input)
			r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
			w := httptest.NewRecorder()

			f := storage.NewFileStorage()
			assert.NoError(t, f.AddNewURL(r.Context(), "https://yandex.ru/", "3LbIJLJ5", ""))
			assert.NoError(t, f.AddNewURL(r.Context(), "https://practicum.ru/", "99XGYq4c", ""))

			userID := uuid.New().String()
			secretKey, _ := hex.DecodeString(common.SecretKey)
			aesBlock, _ := aes.NewCipher(secretKey)
			aesGCM, _ := cipher.NewGCM(aesBlock)

			nonce := make([]byte, aesGCM.NonceSize())
			io.ReadFull(rand.Reader, nonce)

			plaintext := fmt.Sprintf("%s:%s", "userID", userID)
			encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

			r.AddCookie(&http.Cookie{
				Name:  "userID",
				Value: base64.URLEncoding.EncodeToString(encryptedValue),
			})

			NewHandler(f).Handle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			url, _, _ := f.GetOriginalURLByShortURL(r.Context(), test.input[0])
			assert.Equal(t, "", url)
			url, _, _ = f.GetOriginalURLByShortURL(r.Context(), test.input[1])
			assert.Equal(t, "", url)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			os.Remove(*config.FileStoragePath)
		})
	}

	for _, test := range tests {
		t.Run(test.name+" pgstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.input)
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
			w := httptest.NewRecorder()

			userID := uuid.New().String()
			secretKey, _ := hex.DecodeString(common.SecretKey)
			aesBlock, _ := aes.NewCipher(secretKey)
			aesGCM, _ := cipher.NewGCM(aesBlock)

			nonce := make([]byte, aesGCM.NonceSize())
			io.ReadFull(rand.Reader, nonce)

			plaintext := fmt.Sprintf("%s:%s", "userID", userID)
			encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

			request.AddCookie(&http.Cookie{
				Name:  "userID",
				Value: base64.URLEncoding.EncodeToString(encryptedValue),
			})

			s.EXPECT().DeleteUserUrls(request.Context(), test.input, userID).Return(nil).AnyTimes()
			NewHandler(s).Handle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
