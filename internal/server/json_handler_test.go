package server

import (
	"bytes"
	"encoding/json"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricJSONHandler(t *testing.T) {
	type send struct {
		metrics     metrics.Metrics
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
				metrics: metrics.Metrics{
					ID:    "PollCount",
					Type:  metrics.StringCounterType,
					Delta: value,
					Value: nil,
				},
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
				metrics: metrics.Metrics{
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
				metrics: metrics.Metrics{
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
				metrics: metrics.Metrics{
					ID:    "PollCount",
					Type:  metrics.StringCounterType,
					Delta: nil,
					Value: nil,
				},
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
				metrics: metrics.Metrics{
					ID:    "PollCount",
					Type:  metrics.StringCounterType,
					Delta: value,
					Value: nil,
				},
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
				metrics: metrics.Metrics{
					ID:    "PollCount",
					Type:  metrics.StringCounterType,
					Delta: value,
					Value: nil,
				},
				contentType: "application/json",
				method:      http.MethodPut,
			},
			want: want{
				code:     405,
				response: methodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ser := NewServer("some_address", new(mockStorage))

			jsonBytes, errReqJSON := json.Marshal(tt.send.metrics)
			require.NoError(t, errReqJSON)

			req, errReq := http.NewRequest(tt.send.method, "/update", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			ser.router.ServeHTTP(w, req)
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
		metrics     metrics.Metrics
		contentType string
		method      string
	}
	type want struct {
		storageValue metrics.MetricType
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
				metrics: metrics.Metrics{
					ID:   "SomeCounter",
					Type: metrics.StringCounterType,
				},
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				storageValue: metrics.Counter(value),
				code:         200,
				response: metrics.Metrics{
					ID:    "SomeCounter",
					Type:  metrics.StringCounterType,
					Delta: &value,
				},
			},
		},
		{
			name: "Invalid metric type",
			send: send{
				metrics: metrics.Metrics{
					ID:   "SomeCounter",
					Type: "something",
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
				metrics: metrics.Metrics{
					ID:   "",
					Type: metrics.StringCounterType,
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
			name: "No value for that metric",
			send: send{
				metrics: metrics.Metrics{
					ID:   "SomeCounter",
					Type: metrics.StringCounterType,
				},
				contentType: "application/json",
				method:      http.MethodPost,
			},
			want: want{
				storageValue: metrics.Counter(0),
				err:          metrics.ErrNoValue,
				code:         404,
				response:     JSONObj{"message": "No value"},
			},
		},
		{
			name: "HTTP method not allowed",
			send: send{
				metrics: metrics.Metrics{
					ID:   "SomeCounter",
					Type: metrics.StringCounterType,
				},
				contentType: "application/json",
				method:      http.MethodPut,
			},
			want: want{
				storageValue: metrics.Counter(value),
				code:         405,
				response:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := new(mockStorage)
			ms.On("Get", tt.send.metrics.Type, tt.send.metrics.ID).Return(tt.want.storageValue, tt.want.err)

			ser := NewServer("some_address", ms)

			jsonBytes, errJSON := json.Marshal(tt.send.metrics)
			require.NoError(t, errJSON)

			req, errReq := http.NewRequest(tt.send.method, "/value", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", tt.send.contentType)
			require.NoError(t, errReq)

			w := httptest.NewRecorder()
			ser.router.ServeHTTP(w, req)
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
