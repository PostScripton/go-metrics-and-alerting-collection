package client

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/key_management/rsakeys"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing/hmac"
)

// HTTPClient позволяет делать запрос на сервер
type HTTPClient struct {
	baseURI   string
	client    *resty.Client
	key       string
	publicKey *rsa.PublicKey
	realIP    string
}

func newHTTPClient(baseURI string, timeout time.Duration, hashKey string, publicKey *rsa.PublicKey, address string) *HTTPClient {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Splitting host and port")
		return nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Error().Err(err).Str("address", address).Str("host", host).Msg("Looking up ip")
		return nil
	}

	return &HTTPClient{
		baseURI: baseURI,
		client: resty.New().
			SetBaseURL(baseURI).
			SetTimeout(timeout),
		key:       hashKey,
		publicKey: publicKey,
		realIP:    ips[0].String(),
	}
}

func (c *HTTPClient) UpdateMetric(metric metrics.Metrics) error {
	return c.updateMetricViaJSON(metric)
}

// updateMetricViaURI обновляет метрику, передавая информацию в URI
func (c *HTTPClient) updateMetricViaURI(metricType string, name string, value string) error {
	log.Debug().
		Str("type", metricType).
		Str("name", name).
		Str("value", value).
		Msg("Updating metric")

	uri := fmt.Sprintf("/update/%s/%s/%s", metricType, name, value)
	res, err := c.client.R().
		SetHeader("Content-Type", "text/plain").
		SetHeader("X-Real-IP", c.realIP).
		Post(uri)
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

// updateMetricViaJSON обновляет метрику, передавая информацию через POST-запрос в JSON формате
func (c *HTTPClient) updateMetricViaJSON(metric metrics.Metrics) error {
	if c.key != "" {
		metric.Hash = metric.ToHexHash(hmac.NewHmacSigner(), c.key)
	}

	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	encryptedBytes, err := rsakeys.Encrypt(c.publicKey, jsonBytes)
	if err != nil {
		return err
	}

	log.Debug().Interface("metric", metric).Msg("Sending a metric to update")

	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	if _, err = gz.Write(encryptedBytes); err != nil {
		return err
	}
	if err = gz.Close(); err != nil {
		return err
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", c.realIP).
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

// BatchUpdateMetrics обновляет пачку метрик одним POST-запросом в JSON формате
func (c *HTTPClient) BatchUpdateMetrics(collection map[string]metrics.Metrics) error {
	length := len(collection)
	log.Printf("Sending a batch of [%d] metrics", length)

	var newCollection = make([]metrics.Metrics, 0, length)
	for _, m := range collection {
		if c.key != "" {
			m.Hash = m.ToHexHash(hmac.NewHmacSigner(), c.key)
		}
		newCollection = append(newCollection, m)

		log.Debug().Interface("metric", m).Send()
	}

	jsonBytes, err := json.Marshal(newCollection)
	if err != nil {
		return err
	}

	encryptedBytes, err := rsakeys.Encrypt(c.publicKey, jsonBytes)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	if _, err := gz.Write(encryptedBytes); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", c.realIP).
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

// Close закрывает HTTP соединение
func (c *HTTPClient) Close() error {
	return nil
}
