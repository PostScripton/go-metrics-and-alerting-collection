package metrics

import (
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing"
)

const (
	StringCounterType = "counter"
	StringGaugeType   = "gauge"
)

type Metrics struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

var ErrNoValue = errors.New("no value")

func New(t string, name string) *Metrics {
	return &Metrics{
		ID:   name,
		Type: t,
	}
}

func NewCounter(name string, delta int64) *Metrics {
	return &Metrics{
		ID:    name,
		Type:  StringCounterType,
		Delta: &delta,
	}
}

func NewGauge(name string, value float64) *Metrics {
	return &Metrics{
		ID:    name,
		Type:  StringGaugeType,
		Value: &value,
	}
}

func Update(old *Metrics, new *Metrics) {
	old.ID = new.ID
	old.Type = new.Type
	old.Hash = new.Hash

	switch new.Type {
	case StringCounterType:
		var delta int64
		if old.Delta != nil {
			delta = *old.Delta
		}
		if new.Delta != nil {
			delta += *new.Delta
		}
		old.Delta = &delta
	case StringGaugeType:
		old.Value = new.Value
	}
}

func (m *Metrics) Validate() (bool, error) {
	if m.ID == "" {
		return false, errors.New("empty metric id")
	}
	if m.Type != StringCounterType && m.Type != StringGaugeType {
		return false, errors.New("unsupported metric type")
	}

	return true, nil
}

func (m *Metrics) ToHash(signer hashing.Signer, key string) []byte {
	var data string
	switch m.Type {
	case StringCounterType:
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.Type, *m.Delta)
	case StringGaugeType:
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.Type, *m.Value)
	}

	return signer.Hash(data, key)
}

func (m *Metrics) ToHexHash(signer hashing.Signer, key string) string {
	return signer.HashToHex(m.ToHash(signer, key))
}

func (m *Metrics) ValidHash(signer hashing.Signer, hash string, key string) bool {
	sign := m.ToHash(signer, key)

	return signer.ValidHash(sign, hash)
}
