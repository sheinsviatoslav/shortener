package createurl

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/mocks"
	"github.com/sheinsviatoslav/shortener/internal/storage"
)

func TestCreateHandler(t *testing.T) {
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
		url  string
		want want
	}{
		{
			name: "url already exists",
			url:  "https://practicum.yandex.ru/",
			want: want{
				code:        409,
				response:    "http://localhost:8080/99XGYq4c",
				contentType: "text/plain",
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
			want: want{
				code:        201,
				response:    "http://localhost:8080/7IENelKX",
				contentType: "text/plain",
				getShortURLReturn: getShortURLByOriginalURLReturn{
					shortURL: "",
					isExists: false,
					error:    nil,
				},
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

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()

			m := storage.NewMemStorage()
			if err := m.AddNewURL(request.Context(), "https://practicum.yandex.ru/", "99XGYq4c", ""); err != nil {
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
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()

			fs := storage.NewFileStorage()
			if err := fs.AddNewURL(request.Context(), "https://practicum.yandex.ru/", "99XGYq4c", ""); err != nil {
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
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()

			s.EXPECT().GetShortURLByOriginalURL(r.Context(), test.url).Return(
				test.want.getShortURLReturn.shortURL,
				test.want.getShortURLReturn.isExists,
				test.want.getShortURLReturn.error,
			).AnyTimes()
			s.EXPECT().AddNewURL(r.Context(), test.url, gomock.Any(), "").Return(nil).AnyTimes()

			NewHandler(s).Handle(w, r)

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
