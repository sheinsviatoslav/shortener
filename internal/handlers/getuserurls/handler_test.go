package getuserurls

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/sheinsviatoslav/shortener/internal/config"
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
	"github.com/sheinsviatoslav/shortener/internal/storage"
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
		name      string
		output    storage.UserUrls
		outputErr error
		want      want
	}{
		{
			name: "success",
			output: storage.UserUrls{
				{
					OriginalURL: "https://yandex.ru/",
					ShortURL:    "http://localhost:8080/3LbIJLJ5",
				},
			},
			want: want{
				code:        200,
				response:    "[{\"short_url\":\"http://localhost:8080/3LbIJLJ5\",\"original_url\":\"https://yandex.ru/\"}]",
				contentType: "application/json",
			},
		},
		{
			name:   "no content",
			output: storage.UserUrls{},
			want: want{
				code:        204,
				response:    "[]",
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			m := storage.NewMemStorage()
			if len(test.output) > 0 {
				assert.NoError(t, m.AddNewURL(r.Context(), "https://yandex.ru/", "3LbIJLJ5", ""))
			}

			NewHandler(m).Handle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	for _, test := range tests {
		t.Run(test.name+" filestorage", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			f := storage.NewFileStorage()
			if len(test.output) > 0 {
				assert.NoError(t, f.AddNewURL(r.Context(), "https://yandex.ru/", "3LbIJLJ5", ""))
			}

			NewHandler(f).Handle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			os.Remove(*config.FileStoragePath)
		})
	}

	for _, test := range tests {
		t.Run(test.name+" pgstorage", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
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

			s.EXPECT().GetUserUrls(request.Context(), userID).Return(test.output, test.outputErr).AnyTimes()
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
