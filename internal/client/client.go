package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURI string
	client  http.Client
}

func New(baseURI string, timeout time.Duration) *Client {
	return &Client{
		baseURI: baseURI,
		client: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) UpdateMetricJSON(metricType string, name string, value interface{}) {
	fmt.Printf("--- [%s] \"%s\": %s\n", metricType, name, value)

	payload := metrics.Metrics{
		ID:   name,
		Type: metricType,
	}
	switch metricType {
	case metrics.StringCounterType:
		counter := value.(metrics.MetricIntCaster).ToInt64()
		payload.Delta = &counter
	case metrics.StringGaugeType:
		gauge := value.(metrics.MetricFloatCaster).ToFloat64()
		payload.Value = &gauge
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("JSON error: %s\n", err.Error())
		return
	}

	url := fmt.Sprintf("%s/update", c.baseURI)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Printf("Request error: %s\n", err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		fmt.Printf("Send request error: %s\n", err.Error())
		return
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, errBodyReader := io.ReadAll(response.Body)
		if errBodyReader != nil {
			fmt.Printf("Reading response body error: %s\n", err.Error())
			return
		}

		fmt.Printf("Status code: %d\n", response.StatusCode)
		fmt.Printf("Message: %s\n", string(body))
	}
}
