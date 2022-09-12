package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/rs/zerolog/log"
)

func HashMetric(metric *metrics.Metrics, key string) []byte {
	switch metric.Type {
	case metrics.StringCounterType:
		return Hash(fmt.Sprintf("%s:%s:%d", metric.ID, metric.Type, *metric.Delta), key)
	case metrics.StringGaugeType:
		return Hash(fmt.Sprintf("%s:%s:%f", metric.ID, metric.Type, *metric.Value), key)
	default:
		return []byte{}
	}
}

func Hash(data string, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}

func HashToHexMetric(metric *metrics.Metrics, key string) string {
	return HashToHex(HashMetric(metric, key))
}

func HashToHex(hash []byte) string {
	return hex.EncodeToString(hash)
}

func ValidHash(sign []byte, hash string) bool {
	data, err := hex.DecodeString(hash)
	if err != nil {
		log.Warn().Err(err).Msg("Decoding hex")
		return false
	}

	return hmac.Equal(sign, data)
}
