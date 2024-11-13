//package geturl
//
//import (
//	"github.com/go-chi/chi/v5"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//var testStorage = map[string]string{
//	"https://yandex.ru": "/454FcJTrKC",
//}
//
//func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
//	req, err := http.NewRequest(method, ts.URL+path, nil)
//	require.NoError(t, err)
//
//	resp, err := ts.Client().Do(req)
//	require.NoError(t, err)
//	defer resp.Body.Close()
//
//	respBody, err := io.ReadAll(resp.Body)
//	require.NoError(t, err)
//
//	return resp, string(respBody)
//}
//
//func TestRouter(t *testing.T) {
//	r := chi.NewRouter()
//	r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
//		Handler(w, req, testStorage)
//	})
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//	}
//	tests := []struct {
//		name string
//		path string
//		want want
//	}{
//		{
//			name: "success",
//			path: "/454FcJTrKC",
//			want: want{
//				code:        307,
//				response:    "https://yandex.ru",
//				contentType: "text/plain",
//			},
//		},
//		{
//			name: "invalid path",
//			path: "/qqqqqq",
//			want: want{
//				code:        400,
//				response:    "invalid URL path\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "empty path",
//			path: "/",
//			want: want{
//				code:        400,
//				response:    "empty path\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			resp, get := testRequest(t, ts, "GET", test.path)
//			assert.Equal(t, test.want.code, resp.StatusCode)
//			assert.Equal(t, test.want, get)
//		})
//	}
//}

package geturl

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testStorage = map[string]string{
	"https://yandex.ru": "454FcJTrKC",
}

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
			id:   "454FcJTrKC",
			want: want{
				code:        307,
				response:    "https://yandex.ru",
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+test.id, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			Handler(w, r, testStorage)

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
