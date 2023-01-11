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

// Client позволяет делать запрос на сервер
type Client struct {
	baseURI   string
	client    *resty.Client
	key       string
	publicKey *rsa.PublicKey
	realIP    string
}

func NewClient(baseURI string, timeout time.Duration, key string, cryptoKey string, address string) *Client {
	publicKey, err := rsakeys.ImportPublicKeyFromFile(cryptoKey)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get public key from file")
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Splitting host and port")
		return nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Looking up ip")
		return nil
	}

	return &Client{
		baseURI: baseURI,
		client: resty.New().
			SetBaseURL(baseURI).
			SetTimeout(timeout),
		key:       key,
		publicKey: publicKey,
		realIP:    ips[0].String(),
	}
}

// UpdateMetric обновляет метрику, передавая информацию в URI
func (c *Client) UpdateMetric(metricType string, name string, value string) error {
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

// UpdateMetricJSON обновляет метрику, передавая информацию через POST-запрос в JSON формате
func (c *Client) UpdateMetricJSON(metric metrics.Metrics) error {
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

// UpdateMetricsBatchJSON обновляет пачку метрик одним POST-запросом в JSON формате
func (c *Client) UpdateMetricsBatchJSON(collection map[string]metrics.Metrics) error {
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
