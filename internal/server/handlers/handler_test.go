package handlers

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Store(name string, value metrics.MetricType) {
}

func (m *mockStorage) Get(t string, name string) (metrics.MetricType, error) {
	args := m.Called(t, name)
	return args.Get(0).(metrics.MetricType), args.Error(1)
}

func TestUpdateMetricHandler(t *testing.T) {
	type send struct {
		uri         string
		contentType string
		method      string
	}
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name string
		send send
		want want
	}{
		{
			name: "OK",
			send: send{
				uri:         "/update/counter/PollCount/5",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name: "Invalid metric type",
			send: send{
				uri:         "/update/something/PollCount/5",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     501,
				response: "Wrong metric type",
			},
		},
		{
			name: "No metric ID specified",
			send: send{
				uri:         "/update/counter",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
		},
		{
			name: "No metric value specified",
			send: send{
				uri:         "/update/counter/SomeCounter",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
		},
		{
			name: "Invalid metric value",
			send: send{
				uri:         "/update/counter/SomeCounter/none",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     400,
				response: "Invalid metric value",
			},
		},
		//{
		//	name: "Invalid Content-Type",
		//	send: send{
		//		uri:         "/update/counter/PollCount/5",
		//		contentType: "application/json",
		//		method:      http.MethodPost,
		//	},
		//	want: want{
		//		code:     400,
		//		response: "Invalid Content-Type",
		//	},
		//},
		{
			name: "HTTP method not allowed",
			send: send{
				uri:         "/update/counter/PollCount/5",
				contentType: "text/plain",
				method:      http.MethodPut,
			},
			want: want{
				code:     405,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post("/update/{type}/{name}/{value}", UpdateMetricHandler(new(mockStorage)))

			req, errReq := http.NewRequest(tt.send.method, tt.send.uri, nil)
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, errReadBody := io.ReadAll(res.Body)
			require.NoError(t, errReadBody)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	type send struct {
		uri         string
		contentType string
		method      string
	}
	type want struct {
		metricType  string
		metricName  string
		metricValue metrics.MetricType
		err         error
		code        int
		response    string
	}
	tests := []struct {
		name string
		send send
		want want
	}{
		{
			name: "OK",
			send: send{
				uri:         "/value/counter/SomeCounter",
				contentType: "text/plain",
				method:      http.MethodGet,
			},
			want: want{
				metricType:  "counter",
				metricName:  "SomeCounter",
				metricValue: metrics.Counter(5),
				err:         nil,
				code:        200,
				response:    "5",
			},
		},
		{
			name: "Wrong metric type",
			send: send{
				uri:         "/value/qwerty/SomeCounter",
				contentType: "text/plain",
				method:      http.MethodGet,
			},
			want: want{
				metricType:  "qwerty",
				metricName:  "SomeCounter",
				metricValue: metrics.Counter(0),
				err:         nil,
				code:        501,
				response:    "Wrong metric type",
			},
		},
		{
			name: "No metric ID specified",
			send: send{
				uri:         "/value/counter",
				contentType: "text/plain",
				method:      http.MethodGet,
			},
			want: want{
				metricType:  "counter",
				metricName:  "",
				metricValue: metrics.Counter(0),
				err:         nil,
				code:        404,
				response:    "404 page not found\n",
			},
		},
		{
			name: "No value for that metric",
			send: send{
				uri:         "/value/counter/SomeCounter",
				contentType: "text/plain",
				method:      http.MethodGet,
			},
			want: want{
				metricType:  "counter",
				metricName:  "SomeCounter",
				metricValue: metrics.Counter(0),
				err:         metrics.ErrNoValue,
				code:        404,
				response:    "",
			},
		},
		{
			name: "HTTP method not allowed",
			send: send{
				uri:         "/value/counter/SomeCounter",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				metricType:  "counter",
				metricName:  "SomeCounter",
				metricValue: metrics.Counter(0),
				err:         nil,
				code:        405,
				response:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := new(mockStorage)
			ms.On("Get", tt.want.metricType, tt.want.metricName).Return(tt.want.metricValue, tt.want.err)

			router := chi.NewRouter()
			router.Get("/value/{type}/{name}", GetMetricHandler(ms))

			req, errReq := http.NewRequest(tt.send.method, tt.send.uri, nil)
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, errReadBody := io.ReadAll(res.Body)
			require.NoError(t, errReadBody)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
