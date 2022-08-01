package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"io"
	"net/http"
	"os"
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

//func (c *Client) UpdateMetric(metricType string, name string, value string) {
//	fmt.Printf("--- [%s] \"%s\": %s\n", metricType, name, value)
//
//	url := fmt.Sprintf("%s/update/%s/%s/%s", c.baseURI, metricType, name, value)
//	request, err := http.NewRequest(http.MethodPost, url, nil)
//	if err != nil {
//		fmt.Printf("Error: %s\n", err.Error())
//		os.Exit(1)
//	}
//	request.Header.Set("Content-Type", "text/plain")
//
//	response, err := c.client.Do(request)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	if response.StatusCode != http.StatusOK {
//		defer response.Body.Close()
//		body, errBodyReader := io.ReadAll(response.Body)
//		if errBodyReader != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		fmt.Printf("Status code: %d\n", response.StatusCode)
//		fmt.Printf("Message: %s\n", string(body))
//	}
//}

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
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/update", c.baseURI)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, errBodyReader := io.ReadAll(response.Body)
		if errBodyReader != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Status code: %d\n", response.StatusCode)
		fmt.Printf("Message: %s\n", string(body))
	}
}
