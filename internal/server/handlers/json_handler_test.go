package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricJSONHandler(t *testing.T) {
	type send struct {
		metrics     *metrics.Metrics
		contentType string
		method      string
	}
	type want struct {
		code     int
		response JSONObj
	}
	var value = new(int64)
	*value = int64(5)
	tests := []struct {
		name string
		send send
		want want
	}{
		{
			name: "OK",
			send: send{
				metrics:     metrics.NewCounter("PollCount", *value),
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				code:     200,
				response: JSONObj{},
			},
		},
		{
			name: "Invalid metric type",
			send: send{
				metrics: &metrics.Metrics{
					ID:    "PollCount",
					Type:  "something",
					Delta: value,
					Value: nil,
				},
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				code:     501,
				response: JSONObj{"message": "Invalid metric type"},
			},
		},
		{
			name: "No metric ID specified",
			send: send{
				metrics: &metrics.Metrics{
					ID:    "",
					Type:  metrics.StringCounterType,
					Delta: value,
					Value: nil,
				},
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				code:     404,
				response: JSONObj{"message": "No metric ID specified"},
			},
		},
		{
			name: "No metric value specified",
			send: send{
				metrics:     metrics.New(metrics.StringCounterType, "PollCount"),
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				code:     404,
				response: notFoundResponse,
			},
		},
		{
			name: "Invalid Content-Type",
			send: send{
				metrics:     metrics.NewCounter("PollCount", *value),
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     400,
				response: JSONObj{"message": "Invalid Content-Type"},
			},
		},
		{
			name: "HTTP method not allowed",
			send: send{
				metrics:     metrics.NewCounter("PollCount", *value),
				contentType: "application/json",
				method:      http.MethodPut,
			},
			want: want{
				code:     405,
				response: notFoundResponse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Get("/update", UpdateMetricHandler(new(mockStorage)))
			router.NotFound(NotFound)

			jsonBytes, errReqJSON := json.Marshal(*tt.send.metrics)
			require.NoError(t, errReqJSON)

			req, errReq := http.NewRequest(tt.send.method, "/update", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, errReadBody := io.ReadAll(res.Body)
			require.NoError(t, errReadBody)

			if len(resBody) == 0 || res.Header.Get("Content-Type") != "application/json" {
				return
			}

			var jsonRes JSONObj
			errResJSON := json.Unmarshal(resBody, &jsonRes)
			require.NoError(t, errResJSON)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, jsonRes)
		})
	}
}

func TestGetMetricJSONHandler(t *testing.T) {
	type send struct {
		metrics     *metrics.Metrics
		contentType string
		method      string
	}
	type want struct {
		metricReturn *metrics.Metrics
		err          error
		code         int
		response     any
	}
	var value int64 = 5
	tests := []struct {
		name string
		send send
		want want
	}{
		{
			name: "OK",
			send: send{
				metrics:     metrics.New(metrics.StringCounterType, "SomeCounter"),
				contentType: "application/json",
				method:      http.MethodGet,
			},
			want: want{
				metricReturn: metrics.NewCounter("SomeCounter", 5),
				code:         200,
				response:     metrics.NewCounter("SomeCounter", value),
			},
		},
		{
			name: "Invalid metric type",
			send: send{
				metrics:     metrics.New("something", "SomeCounter"),
				contentType: "application/json",
				method:      http.MethodGet,
			},
			want: want{
				code:     501,
				response: JSONObj{"message": "Invalid metric type"},
			},
		},
		{
			name: "No metric ID specified",
			send: send{
				metrics:     metrics.New(metrics.StringCounterType, ""),
				contentType: "application/json",
				method:      http.MethodGet,
			},
			want: want{
				code:     404,
				response: JSONObj{"message": "No metric ID specified"},
			},
		},
		{
			name: "No value for that metric",
			send: send{
				metrics:     metrics.New(metrics.StringCounterType, "SomeCounter"),
				contentType: "application/json",
				method:      http.MethodGet,
			},
			want: want{
				err:      metrics.ErrNoValue,
				code:     404,
				response: JSONObj{"message": "No value"},
			},
		},
		{
			name: "HTTP method not allowed",
			send: send{
				metrics:     metrics.New(metrics.StringCounterType, "SomeCounter"),
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				code:     405,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := new(mockStorage)
			ms.On("Get", *metrics.New(tt.send.metrics.Type, tt.send.metrics.ID)).Return(tt.want.metricReturn, tt.want.err)

			router := chi.NewRouter()
			router.Get("/value", GetMetricJSONHandler(ms))

			jsonBytes, errJSON := json.Marshal(tt.send.metrics)
			require.NoError(t, errJSON)

			req, errReq := http.NewRequest(tt.send.method, "/value", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, errReadBody := io.ReadAll(res.Body)
			require.NoError(t, errReadBody)

			if len(resBody) == 0 || res.Header.Get("Content-Type") != "application/json" {
				return
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
			switch tt.want.response.(type) {
			case metrics.Metrics:
				var jsonRes metrics.Metrics
				errResJSON := json.Unmarshal(resBody, &jsonRes)
				require.NoError(t, errResJSON)
				assert.Equal(t, tt.want.response, jsonRes)
			case JSONObj:
				var jsonRes JSONObj
				errResJSON := json.Unmarshal(resBody, &jsonRes)
				require.NoError(t, errResJSON)
				assert.Equal(t, tt.want.response, jsonRes)
			}
		})
	}
}
