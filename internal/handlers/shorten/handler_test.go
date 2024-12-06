package shorten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/mocks"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStorage(ctrl)

	type getShortURLByOriginalURLReturn struct {
		shortURL string
		isExists bool
		error    error
	}

	type want struct {
		code              int
		response          string
		contentType       string
		getShortURLReturn getShortURLByOriginalURLReturn
	}
	tests := []struct {
		name string
		body map[string]interface{}
		url  string
		want want
	}{
		{
			name: "url already exists",
			url:  "https://practicum.yandex.ru/",
			body: map[string]interface{}{
				"url": "https://practicum.yandex.ru/",
			},
			want: want{
				code:        409,
				response:    `{"result":"http://localhost:8080/99XGYq4c"}`,
				contentType: "application/json",
				getShortURLReturn: getShortURLByOriginalURLReturn{
					shortURL: "99XGYq4c",
					isExists: true,
					error:    nil,
				},
			},
		},
		{
			name: "successfully created",
			url:  "https://yandex.ru/",
			body: map[string]interface{}{
				"url": "https://yandex.ru/",
			},
			want: want{
				code:        201,
				response:    `{"result":"http://localhost:8080/7IENelKX"}`,
				contentType: "application/json",
				getShortURLReturn: getShortURLByOriginalURLReturn{
					shortURL: "",
					isExists: false,
					error:    nil,
				},
			},
		},
		{
			name: "invalid url",
			url:  "",
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
			url:  "",
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
			url:  "https://yandex.ru",
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
			url:  "",
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

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()

			m := storage.NewMemStorage()
			if err := m.AddNewURL("https://practicum.yandex.ru/", "99XGYq4c"); err != nil {
				require.NoError(t, err)
			}
			NewHandler(m).Handle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			if test.name == "successfully created" {
				isMatch, _ := regexp.MatchString(fmt.Sprintf(
					"http://localhost:8080/[0-9a-zA-Z]{%d}",
					common.DefaultHashLength,
				), string(resBody))
				require.NoError(t, err)
				assert.Equal(t, true, isMatch)
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	for _, test := range tests {
		t.Run(test.name+" filestorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()

			fs := storage.NewFileStorage()
			if err := fs.AddNewURL("https://practicum.yandex.ru/", "99XGYq4c"); err != nil {
				require.NoError(t, err)
			}
			NewHandler(fs).Handle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			if test.name == "successfully created" {
				isMatch, _ := regexp.MatchString(fmt.Sprintf(
					"http://localhost:8080/[0-9a-zA-Z]{%d}",
					common.DefaultHashLength,
				), string(resBody))
				require.NoError(t, err)
				assert.Equal(t, true, isMatch)
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			err = os.Remove(*config.FileStoragePath)
			if err != nil {
				require.NoError(t, err)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name+" pgstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()

			s.EXPECT().GetShortURLByOriginalURL(test.url).Return(
				test.want.getShortURLReturn.shortURL,
				test.want.getShortURLReturn.isExists,
				test.want.getShortURLReturn.error,
			).AnyTimes()
			s.EXPECT().AddNewURL(test.url, gomock.Any()).Return(nil).AnyTimes()

			NewHandler(s).Handle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			if test.name == "successfully created" {
				isMatch, _ := regexp.MatchString(fmt.Sprintf(
					"http://localhost:8080/[0-9a-zA-Z]{%d}",
					common.DefaultHashLength,
				), string(resBody))
				require.NoError(t, err)
				assert.Equal(t, true, isMatch)
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
