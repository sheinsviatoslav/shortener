package shortenbatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sheinsviatoslav/shortener/internal/mocks"
	"github.com/sheinsviatoslav/shortener/internal/storage"
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
		name      string
		body      storage.InputManyUrls
		output    storage.OutputManyUrls
		outputErr error
		want      want
	}{
		{
			name: "successfully created",
			body: storage.InputManyUrls{
				{
					CorrelationID: "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
					OriginalURL:   "https://practicum.yandex.ru/",
				},
				{
					CorrelationID: "489748a4-7521-4821-bdf9-52c1f6059387",
					OriginalURL:   "https://yandex.ru/",
				},
			},
			output: storage.OutputManyUrls{
				{
					CorrelationID: "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
					ShortURL:      "http://localhost:8080/3LbIJLJ5",
				},
				{
					CorrelationID: "489748a4-7521-4821-bdf9-52c1f6059387",
					ShortURL:      "http://localhost:8080/7wcVQIE1",
				},
			},
			want: want{
				code: 201,
				response: "[{\"correlation_id\":\"0a7b6ee4-ffdf-4394-9ac6-b42b652b389a\",\"short_url\":\"http://localhost:8080/3LbIJLJ5\"}," +
					"{\"correlation_id\":\"489748a4-7521-4821-bdf9-52c1f6059387\",\"short_url\":\"http://localhost:8080/7wcVQIE1\"}]",
				contentType: "application/json",
			},
		},
		{
			name:   "empty body",
			body:   storage.InputManyUrls{},
			output: storage.OutputManyUrls{},
			want: want{
				code:        400,
				response:    "empty request body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		//{
		//	name: "output err",
		//	body: storage.InputManyUrls{
		//		{
		//			CorrelationID: "0a7b6ee4-ffdf-4394-9ac6-b42b652b389a",
		//			OriginalURL:   "https://practicum.yandex.ru/",
		//		},
		//		{
		//			CorrelationID: "489748a4-7521-4821-bdf9-52c1f6059387",
		//			OriginalURL:   "https://yandex.ru/",
		//		},
		//	},
		//	output:    nil,
		//	outputErr: errors.New("output error"),
		//	want: want{
		//		code:        400,
		//		response:    "output error\n",
		//		contentType: "text/plain; charset=utf-8",
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name+" memstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
			w := httptest.NewRecorder()

			m := storage.NewMemStorage()
			NewHandler(m).Handle(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			var resBody storage.OutputManyUrls
			var buf bytes.Buffer
			buf.ReadFrom(res.Body)
			json.Unmarshal(buf.Bytes(), &resBody)

			for _, v := range resBody {
				isMatch, _ := regexp.MatchString(fmt.Sprintf(
					"http://localhost:8080/[0-9a-zA-Z]{%d}",
					common.DefaultHashLength,
				), v.ShortURL)
				assert.Equal(t, true, isMatch)
			}

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	for _, test := range tests {
		t.Run(test.name+" filestorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
			w := httptest.NewRecorder()

			fs := storage.NewFileStorage()
			NewHandler(fs).Handle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			var resBody storage.OutputManyUrls
			var buf bytes.Buffer
			buf.ReadFrom(res.Body)
			json.Unmarshal(buf.Bytes(), &resBody)

			for _, v := range resBody {
				isMatch, _ := regexp.MatchString(fmt.Sprintf(
					"http://localhost:8080/[0-9a-zA-Z]{%d}",
					common.DefaultHashLength,
				), v.ShortURL)
				assert.Equal(t, true, isMatch)
			}

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			os.Remove(*config.FileStoragePath)
		})
	}

	for _, test := range tests {
		t.Run(test.name+" pgstorage", func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
			w := httptest.NewRecorder()

			s.EXPECT().AddManyUrls(request.Context(), test.body, "").Return(test.output, test.outputErr).AnyTimes()
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
