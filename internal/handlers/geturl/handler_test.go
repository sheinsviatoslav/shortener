package geturl

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/mocks"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStorage(ctrl)

	type getOriginalURLByShortURLReturn struct {
		originalURL string
		error       error
	}

	type want struct {
		code                 int
		response             string
		contentType          string
		getOriginalURLReturn getOriginalURLByShortURLReturn
	}
	tests := []struct {
		name     string
		shortURL string
		want     want
	}{
		{
			name:     "success",
			shortURL: "99XGYq4c",
			want: want{
				code:        307,
				response:    "https://practicum.yandex.ru/",
				contentType: "text/plain",
				getOriginalURLReturn: getOriginalURLByShortURLReturn{
					originalURL: "https://practicum.yandex.ru/",
					error:       nil,
				},
			},
		},
		{
			name:     "invalid path",
			shortURL: "qqqqqq",
			want: want{
				code:        400,
				response:    "invalid URL path\n",
				contentType: "text/plain; charset=utf-8",
				getOriginalURLReturn: getOriginalURLByShortURLReturn{
					originalURL: "",
					error:       errors.New("invalid URL path"),
				},
			},
		},
		{
			name:     "empty path",
			shortURL: "",
			want: want{
				code:        400,
				response:    "empty path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+test.shortURL, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortURL", test.shortURL)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			m := storage.NewMemStorage()
			if err := m.AddNewURL("https://practicum.yandex.ru/", "99XGYq4c", ""); err != nil {
				require.NoError(t, err)
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
			r := httptest.NewRequest(http.MethodGet, "/"+test.shortURL, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortURL", test.shortURL)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			fs := storage.NewFileStorage()
			if err := fs.AddNewURL("https://practicum.yandex.ru/", "99XGYq4c", ""); err != nil {
				require.NoError(t, err)
			}
			NewHandler(fs).Handle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			err = os.Remove(*config.FileStoragePath)
			if err != nil {
				require.NoError(t, err)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name+" pgstorage", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+test.shortURL, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortURL", test.shortURL)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			s.EXPECT().GetOriginalURLByShortURL(test.shortURL).Return(
				test.want.getOriginalURLReturn.originalURL,
				test.want.getOriginalURLReturn.error,
			).AnyTimes()
			NewHandler(s).Handle(w, r)

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
