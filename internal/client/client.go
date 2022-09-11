package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/hashing"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURI string
	client  *resty.Client
	key     string
}

func NewClient(baseURI string, timeout time.Duration, key string) *Client {
	return &Client{
		baseURI: baseURI,
		client: resty.New().
			SetBaseURL(baseURI).
			SetRetryCount(5).
			SetRetryWaitTime(10 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			SetTimeout(timeout),
		key: key,
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

func (c *Client) UpdateMetricJSON(metric metrics.Metrics) error {
	if c.key != "" {
		metric.Hash = hashing.HashToHexMetric(&metric, c.key)
	}

	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	fmt.Printf("--- Send [JSON] [%s]\n", string(jsonBytes))

	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	if _, err := gz.Write(jsonBytes); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(out.Bytes()).
		Post("/update")
	if err != nil {
		return fmt.Errorf("send request error: %s", err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("--- Response ---\nStatus code: [%d]\nMessage: [%s]", res.StatusCode(), strings.Trim(string(res.Body()), "\n"))
	}

	return nil
}

func (c *Client) UpdateMetricsBatchJSON(collection map[string]metrics.Metrics) error {
	var newCollection = make([]metrics.Metrics, 0, len(collection))
	for _, m := range collection {
		if c.key != "" {
			m.Hash = hashing.HashToHexMetric(&m, c.key)
		}
		newCollection = append(newCollection, m)
	}

	jsonBytes, err := json.Marshal(newCollection)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	if _, err := gz.Write(jsonBytes); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	fmt.Printf("Sending a batch of [%d] metrics\n", len(newCollection))

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(out.Bytes()).
		Post("/updates")
	if err != nil {
		return fmt.Errorf("send request error: %s", err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("--- Response ---\nStatus code: [%d]\nMessage: [%s]", res.StatusCode(), strings.Trim(string(res.Body()), "\n"))
	}

	return nil
}
