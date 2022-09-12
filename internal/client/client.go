package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/hashing"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
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
			SetTimeout(timeout),
		key: key,
	}
}

func (c *Client) UpdateMetric(metricType string, name string, value string) error {
	log.Debug().
		Str("type", metricType).
		Str("name", name).
		Str("value", value).
		Msg("Updating metric")

	uri := fmt.Sprintf("/update/%s/%s/%s", metricType, name, value)
	res, err := c.client.R().SetHeader("Content-Type", "text/plain").Post(uri)
	if err != nil {
		return fmt.Errorf("send request error: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		message := strings.Trim(string(res.Body()), "\n")
		log.Warn().Int("status_code", res.StatusCode()).Str("message", message).Msg("Response")
		return fmt.Errorf(message)
	}

	return nil
}

func (c *Client) UpdateMetricJSON(metric metrics.Metrics) error {
	if c.key != "" {
		metric.Hash = hashing.HashToHexMetric(&metric, c.key)
	}

	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	log.Debug().Interface("metric", metric).Msg("Sending a metric to update")

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
		return fmt.Errorf("send request error: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		message := strings.Trim(string(res.Body()), "\n")
		log.Warn().Int("status_code", res.StatusCode()).Str("message", message).Msg("Response")
		return fmt.Errorf(message)
	}

	return nil
}

func (c *Client) UpdateMetricsBatchJSON(collection map[string]metrics.Metrics) error {
	length := len(collection)
	log.Printf("Sending a batch of [%d] metrics", length)

	var newCollection = make([]metrics.Metrics, 0, length)
	for _, m := range collection {
		if c.key != "" {
			m.Hash = hashing.HashToHexMetric(&m, c.key)
		}
		newCollection = append(newCollection, m)

		log.Debug().Interface("metric", m).Send()
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

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(out.Bytes()).
		Post("/updates")
	if err != nil {
		return fmt.Errorf("send request error: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		message := strings.Trim(string(res.Body()), "\n")
		log.Warn().Int("status_code", res.StatusCode()).Str("message", message).Msg("Response")
		return fmt.Errorf(message)
	}

	return nil
}
