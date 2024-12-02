package shorten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/handlers/createurl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
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
		body map[string]interface{}
		want want
	}{
		{
			name: "url already exists",
			body: map[string]interface{}{
				"url": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        200,
				response:    "{\"result\":\"http://localhost:8080/99XGYq4c\"}",
				contentType: "application/json",
			},
		},
		{
			name: "invalid url",
			body: map[string]interface{}{
				"url": "h",
			},
			want: want{
				code:        400,
				response:    "parse \"h\": invalid URI for request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "unable to parse input json",
			body: map[string]interface{}{
				"url": make(chan int),
			},
			want: want{
				code:        400,
				response:    "unexpected end of JSON input\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "no url param",
			body: map[string]interface{}{
				"myURL": "https://yandex.ru",
			},
			want: want{
				code:        400,
				response:    "url is required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "empty url",
			body: map[string]interface{}{
				"url": "",
			},
			want: want{
				code:        400,
				response:    "url is required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	t.Setenv("FILE_STORAGE_PATH", "mocks/url_storage_already_exists.json")
	config.Init()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()
			Handler(w, request)

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
		fileName := "mocks/url_storage_create_new_item.json"
		t.Setenv("FILE_STORAGE_PATH", fileName)
		config.Init()
		body, _ := json.Marshal(map[string]interface{}{
			"url": "https://practicum.yandex.ru/",
		})

		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()
		Handler(w, request)

		res := w.Result()
		assert.Equal(t, 201, res.StatusCode)
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		isMatch, _ := regexp.MatchString(fmt.Sprintf("http://localhost:8080/[0-9a-zA-Z]{%d}", createurl.DefaultHashLength), string(resBody))
		require.NoError(t, err)

		assert.Equal(t, true, isMatch)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		e := os.Remove(fileName)
		if e != nil {
			require.NoError(t, err)
		}
	})

}
