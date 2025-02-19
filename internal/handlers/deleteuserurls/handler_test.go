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
	"io"
	"net/http"
	"net/http/httptest"
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
