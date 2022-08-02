package client

import (
	"encoding/json"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type Client struct {
	baseURI string
	client  *resty.Client
}

func New(baseURI string, timeout time.Duration) *Client {
	return &Client{
		baseURI: baseURI,
		client: resty.New().
			SetBaseURL(baseURI).
			SetRetryCount(5).
			SetRetryWaitTime(10 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			SetTimeout(timeout),
	}
}

func (c *Client) UpdateMetric(metricType string, name string, value string) {
	fmt.Printf("--- [%s] \"%s\": %s\n", metricType, name, value)

	uri := fmt.Sprintf("/update/%s/%s/%s", metricType, name, value)
	res, err := c.client.R().SetHeader("Content-Type", "text/plain").Post(uri)
	if err != nil {
		fmt.Printf("Send request error: %s\n", err.Error())
		return
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Printf("Status code: %d\n", res.StatusCode())
		fmt.Printf("Message: %s\n", string(res.Body()))
	}
}

func (c *Client) UpdateMetricJSON(metricType string, name string, value interface{}) {
	fmt.Printf("--- [JSON] [%s] \"%s\": %s\n", metricType, name, value)

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
	jsonBytes, errMarshal := json.Marshal(payload)
	if errMarshal != nil {
		fmt.Printf("JSON error: %s\n", errMarshal.Error())
		return
	}

	res, err := c.client.R().SetHeader("Content-Type", "application/json").SetBody(jsonBytes).Post("/update")
	if err != nil {
		fmt.Printf("Send request error: %s\n", err.Error())
		return
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Printf("Status code: %d\n", res.StatusCode())
		fmt.Printf("Message: %s\n", string(res.Body()))
	}
}
