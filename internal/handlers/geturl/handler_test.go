package geturl

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sheinsviatoslav/shortener/internal/config"
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
		id   string
		want want
	}{
		{
			name: "success",
			id:   "99XGYq4c",
			want: want{
				code:        307,
				response:    "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
		},
		{
			name: "invalid path",
			id:   "qqqqqq",
			want: want{
				code:        400,
				response:    "invalid URL path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "empty path",
			id:   "",
			want: want{
				code:        400,
				response:    "empty path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	t.Setenv("FILE_STORAGE_PATH", "mocks/url_storage.json")
	config.Init()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, "/"+test.id, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			Handler(w, r)

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
