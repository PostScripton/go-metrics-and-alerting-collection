package handlers

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStorage struct{}

func (m mockStorage) Store(name string, value metrics.MetricType) {
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
			name: "Invalid metric type",
			send: send{
				uri:         "/update/counter",
				contentType: "text/plain",
				method:      http.MethodPost,
			},
			want: want{
				code:     404,
				response: "No metric ID specified",
			},
		},
		{
			name: "Invalid metric value",
			send: send{
				uri:         "/update/counter/none",
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
				response: "Method not allowed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.send.method, tt.send.uri, nil)
			req.Header.Set("Content-Type", tt.send.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateMetricHandler(service.NewMetricService(&mockStorage{})))
			h.ServeHTTP(w, req)
			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
