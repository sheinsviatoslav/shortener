package geturl

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "success",
			path: "/454FcJTrKC",
			want: want{
				code:        307,
				response:    "https://yandex.ru",
				contentType: "text/plain",
			},
		},
		{
			name: "invalid path",
			path: "/qqqqqq",
			want: want{
				code:        400,
				response:    "invalid URL path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "empty path",
			path: "/",
			want: want{
				code:        400,
				response:    "empty path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	storage := map[string]string{
		"https://yandex.ru": "/454FcJTrKC",
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			w := httptest.NewRecorder()
			GetHandler(w, request, storage)

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
