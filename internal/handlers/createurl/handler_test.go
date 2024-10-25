package createurl

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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
		url  string
		want want
	}{
		{
			name: "url already exists",
			url:  "https://yandex.ru",
			want: want{
				code:        200,
				response:    "http://example.com/454FcJTrKC",
				contentType: "text/plain",
			},
		},
		{
			name: "invalid url",
			url:  "h",
			want: want{
				code:        400,
				response:    "parse \"h\": invalid URI for request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "empty url",
			url:  "",
			want: want{
				code:        400,
				response:    "url is required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	storage := map[string]string{
		"https://yandex.ru": "/454FcJTrKC",
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()
			PostHandler(w, request, storage)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	t.Run("successfully created", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/"))
		w := httptest.NewRecorder()
		PostHandler(w, request, storage)

		res := w.Result()
		assert.Equal(t, 201, res.StatusCode)
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		isMatch, _ := regexp.MatchString("http://example.com/[0-9a-zA-Z]{8}", string(resBody))
		require.NoError(t, err)

		assert.Equal(t, true, isMatch)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	})

}
