package metrics

import (
	"errors"
	"fmt"

	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing"
)

// Типы метрик
const (
	StringCounterType = "counter" // Счётчик
	StringGaugeType   = "gauge"   // Измеритель
)

// Metrics объект для хранения информации о метрике
type Metrics struct {
	ID    string   `json:"id"`              // Название метрики
	Type  string   `json:"type"`            // Тип метрики
	Delta *int64   `json:"delta,omitempty"` // (counter) Дельта, на которую изменилась метрика
	Value *float64 `json:"value,omitempty"` // (gauge) Новое значение метрики
	Hash  string   `json:"hash,omitempty"`  // Захэшированное значение метрики с помощью HMAC и сохранённое в hex
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

// Update позволяет обновить старую метрику значениями из новой метрики
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

// Validate проверяет корректность полей метрики
func (m *Metrics) Validate() (bool, error) {
	if m.ID == "" {
		return false, errors.New("empty metric id")
	}
	if m.Type != StringCounterType && m.Type != StringGaugeType {
		return false, errors.New("unsupported metric type")
	}

	return true, nil
}

// ToHash хэширует метрику с помощью hashing.Signer
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

// ToHexHash хэширует метрику с помощью hashing.Signer и переводит в hex
func (m *Metrics) ToHexHash(signer hashing.Signer, key string) string {
	return signer.HashToHex(m.ToHash(signer, key))
}

// ValidHash проверяет метрику на подлинность с помощью hashing.Signer
func (m *Metrics) ValidHash(signer hashing.Signer, hash string, key string) bool {
	sign := m.ToHash(signer, key)

	return signer.ValidHash(sign, hash)
}
