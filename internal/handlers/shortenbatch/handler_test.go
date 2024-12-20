package shortenbatch

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/sheinsviatoslav/shortener/internal/common"
//	"github.com/sheinsviatoslav/shortener/internal/config"
//	"github.com/sheinsviatoslav/shortener/internal/storage"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"regexp"
//	"testing"
//)
//
//func TestShortenHandler(t *testing.T) {
//	//ctrl := gomock.NewController(t)
//	//defer ctrl.Finish()
//	//
//	//s := mocks.NewMockStorage(ctrl)
//
//	type getShortURLByOriginalURLReturn struct {
//		shortURL string
//		isExists bool
//		error    error
//	}
//
//	type want struct {
//		code              int
//		response          string
//		contentType       string
//		getShortURLReturn getShortURLByOriginalURLReturn
//	}
//	tests := []struct {
//		name string
//		body []map[string]interface{}
//		want want
//	}{
//		{
//			name: "successfully created",
//			body: []map[string]interface{}{
//				{
//					"correlation_id": "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					"original_url":   "https://practicum.yandex.ru/",
//				},
//				{
//					"correlation_id": "489748a4-7521-4821-bdf9-52c1f6059387",
//					"original_url":   "https://yandex.ru/",
//				},
//			},
//			want: want{
//				code: 201,
//				response: "[{\"correlation_id\":\"0a7b6ee4-ffdf-4394-9ac6-b42b652b389a\",\"short_url\":\"http://localhost:8080/3LbIJLJ5\"}," +
//					"{\"correlation_id\":\"489748a4-7521-4821-bdf9-52c1f6059387\",\"short_url\":\"http://localhost:8080/7wcVQIE1\"}]",
//				contentType: "application/json",
//				getShortURLReturn: getShortURLByOriginalURLReturn{
//					shortURL: "99XGYq4c",
//					isExists: true,
//					error:    nil,
//				},
//			},
//		},
//		{
//			name: "invalid url",
//			body: []map[string]interface{}{
//				{
//					"correlation_id": "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					"original_url":   "h",
//				},
//				{
//					"correlation_id": "489748a4-7521-4821-bdf9-52c1f6059387",
//					"original_url":   "https://yandex.ru/",
//				},
//			},
//			want: want{
//				code:        400,
//				response:    "parse \"h\": invalid URI for request\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "unable to parse input json",
//			body: []map[string]interface{}{
//				{
//					"correlation_id": "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					"original_url":   make(chan int),
//				},
//				{
//					"correlation_id": "489748a4-7521-4821-bdf9-52c1f6059387",
//					"original_url":   "https://yandex.ru/",
//				},
//			},
//			want: want{
//				code:        400,
//				response:    "unexpected end of JSON input\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "no url param",
//			body: []map[string]interface{}{
//				{
//					"correlation_id": "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					"my_url":         "https://practicum.yandex.ru/",
//				},
//				{
//					"correlation_id": "489748a4-7521-4821-bdf9-52c1f6059387",
//					"my_url":         "https://yandex.ru/",
//				},
//			},
//			want: want{
//				code:        400,
//				response:    "url is required\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "empty url",
//			body: []map[string]interface{}{
//				{
//					"correlation_id": "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					"original_url":   "https://practicum.yandex.ru/",
//				},
//				{
//					"correlation_id": "489748a4-7521-4821-bdf9-52c1f6059387",
//					"original_url":   "",
//				},
//			},
//			want: want{
//				code:        400,
//				response:    "url is required\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name+" memstorage", func(t *testing.T) {
//			body, _ := json.Marshal(test.body)
//			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
//			w := httptest.NewRecorder()
//
//			m := storage.NewMemStorage()
//			_, err := m.AddManyUrls(storage.InputManyUrls{
//				{
//					CorrelationID: "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					OriginalURL:   "https://practicum.yandex.ru/",
//				},
//				{
//					CorrelationID: "489748a4-7521-4821-bdf9-52c1f6059387",
//					OriginalURL:   "https://yandex.ru/",
//				},
//			})
//			if err != nil {
//				require.NoError(t, err)
//			}
//			NewHandler(m).Handle(w, request)
//
//			res := w.Result()
//			assert.Equal(t, test.want.code, res.StatusCode)
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//
//			isMatch, _ := regexp.MatchString(fmt.Sprintf(
//				"http://localhost:8080/[0-9a-zA-Z]{%d}",
//				common.DefaultHashLength,
//			), string(resBody))
//			require.NoError(t, err)
//			assert.Equal(t, true, isMatch)
//			assert.Equal(t, test.want.response, string(resBody))
//			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//		})
//	}
//
//	for _, test := range tests {
//		t.Run(test.name+" filestorage", func(t *testing.T) {
//			body, _ := json.Marshal(test.body)
//			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
//			w := httptest.NewRecorder()
//
//			fs := storage.NewFileStorage()
//			_, err := fs.AddManyUrls(storage.InputManyUrls{
//				{
//					CorrelationID: "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
//					OriginalURL:   "https://practicum.yandex.ru/",
//				},
//				{
//					CorrelationID: "489748a4-7521-4821-bdf9-52c1f6059387",
//					OriginalURL:   "https://yandex.ru/",
//				},
//			})
//
//			if err != nil {
//				require.NoError(t, err)
//			}
//			NewHandler(fs).Handle(w, request)
//
//			res := w.Result()
//			assert.Equal(t, test.want.code, res.StatusCode)
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//
//			for _, val := range resBody {
//				fmt.Println(val)
//				isMatch, _ := regexp.MatchString(fmt.Sprintf(
//					"http://localhost:8080/[0-9a-zA-Z]{%d}",
//					common.DefaultHashLength,
//				), string(val))
//				require.NoError(t, err)
//				assert.Equal(t, true, isMatch)
//			}
//
//			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//
//			err = os.Remove(*config.FileStoragePath)
//			if err != nil {
//				require.NoError(t, err)
//			}
//		})
//	}
//
//	//for _, test := range tests {
//	//	t.Run(test.name+" pgstorage", func(t *testing.T) {
//	//		body, _ := json.Marshal(test.body)
//	//		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
//	//		w := httptest.NewRecorder()
//	//
//	//		s.EXPECT().GetShortURLByOriginalURL(test.url).Return(
//	//			test.want.getShortURLReturn.shortURL,
//	//			test.want.getShortURLReturn.isExists,
//	//			test.want.getShortURLReturn.error,
//	//		).AnyTimes()
//	//		s.EXPECT().AddManyUrls(body).Return(nil, nil).AnyTimes()
//	//
//	//		NewHandler(s).Handle(w, request)
//	//
//	//		res := w.Result()
//	//		assert.Equal(t, test.want.code, res.StatusCode)
//	//		defer res.Body.Close()
//	//		resBody, err := io.ReadAll(res.Body)
//	//
//	//		require.NoError(t, err)
//	//
//	//		if test.name == "successfully created" {
//	//			isMatch, _ := regexp.MatchString(fmt.Sprintf(
//	//				"http://localhost:8080/[0-9a-zA-Z]{%d}",
//	//				common.DefaultHashLength,
//	//			), string(resBody))
//	//			require.NoError(t, err)
//	//			assert.Equal(t, true, isMatch)
//	//		} else {
//	//			assert.Equal(t, test.want.response, string(resBody))
//	//		}
//	//		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//	//	})
//	//}
//}
