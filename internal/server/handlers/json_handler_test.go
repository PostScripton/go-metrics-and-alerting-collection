package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/gin-gonic/gin"
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
		response gin.H
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
				response: gin.H{},
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
				response: gin.H{"message": "Invalid metric type"},
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
				response: gin.H{"message": "No metric ID specified"},
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
				response: gin.H{"message": "Invalid Content-Type"},
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
				code:     404,
				response: notFoundResponse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.RedirectTrailingSlash = false
			router.POST("/update", UpdateMetricHandler(new(mockStorage)))
			router.NoRoute(NoRoute)

			jsonBytes, errReqJSON := json.Marshal(tt.send.metrics)
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

			var jsonRes gin.H
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
		response     interface{}
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
				method:      http.MethodGet,
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
				method:      http.MethodGet,
			},
			want: want{
				code:     501,
				response: gin.H{"message": "Invalid metric type"},
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
				method:      http.MethodGet,
			},
			want: want{
				code:     404,
				response: gin.H{"message": "No metric ID specified"},
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
				method:      http.MethodGet,
			},
			want: want{
				storageValue: metrics.Counter(0),
				err:          metrics.ErrNoValue,
				code:         404,
				response:     notFoundResponse,
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
				method:      http.MethodPost,
			},
			want: want{
				storageValue: metrics.Counter(value),
				code:         404,
				response:     notFoundResponse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := new(mockStorage)
			ms.On("Get", tt.send.metrics.Type, tt.send.metrics.ID).Return(tt.want.storageValue, tt.want.err)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.GET("/value", GetMetricHandler(ms))

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

			var jsonRes interface{}
			errResJSON := json.Unmarshal(resBody, &jsonRes)
			require.NoError(t, errResJSON)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, jsonRes)
		})
	}
}
