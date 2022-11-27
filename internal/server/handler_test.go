package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Store(metric metrics.Metrics) error {
	return nil
}

func (m *mockStorage) Get(metric metrics.Metrics) (*metrics.Metrics, error) {
	args := m.Called(metric)
	return args.Get(0).(*metrics.Metrics), args.Error(1)
}

func (m *mockStorage) GetCollection() (map[string]metrics.Metrics, error) {
	args := m.Called()
	return args.Get(0).(map[string]metrics.Metrics), args.Error(1)
}

func (m *mockStorage) StoreCollection(metricsCollection map[string]metrics.Metrics) error {
	args := m.Called(metricsCollection)
	return args.Error(0)
}

func (m *mockStorage) CleanUp() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockStorage) Close() {
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
				response: "405 method not allowed\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ser := NewServer("some_address", new(mockStorage), "")

			req, errReq := http.NewRequest(tt.send.method, tt.send.uri, nil)
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			ser.router.ServeHTTP(w, req)
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
		metricGet    *metrics.Metrics
		metricReturn *metrics.Metrics
		err          error
		code         int
		response     string
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
				metricGet:    metrics.New(metrics.StringCounterType, "SomeCounter"),
				metricReturn: metrics.NewCounter("SomeCounter", 5),
				err:          nil,
				code:         200,
				response:     "5",
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
				metricGet: metrics.New(metrics.StringCounterType, "SomeCounter"),
				err:       nil,
				code:      501,
				response:  "Wrong metric type",
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
				metricGet: metrics.New(metrics.StringCounterType, "SomeCounter"),
				err:       nil,
				code:      404,
				response:  "404 page not found\n",
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
				metricGet:    metrics.New(metrics.StringCounterType, "SomeCounter"),
				metricReturn: nil,
				err:          metrics.ErrNoValue,
				code:         404,
				response:     "",
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
				metricGet: metrics.New(metrics.StringCounterType, "SomeCounter"),
				err:       nil,
				code:      405,
				response:  "405 method not allowed\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := new(mockStorage)
			ms.On("Get", *metrics.New(tt.want.metricGet.Type, tt.want.metricGet.ID)).Return(tt.want.metricReturn, tt.want.err)

			ser := NewServer("some_address", ms, "")

			req, errReq := http.NewRequest(tt.send.method, tt.send.uri, nil)
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			ser.router.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, errReadBody := io.ReadAll(res.Body)
			require.NoError(t, errReadBody)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
