package metrics

import (
	"errors"
	"strconv"
)

const (
	StringCounterType = "counter"
	StringGaugeType   = "gauge"
)

type MetricType interface {
	Type() string
}

type MetricIntCaster interface {
	ToInt64() int64
}

type MetricFloatCaster interface {
	ToFloat64() float64
}

type Metrics struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Counter int64

type Gauge float64

var ErrNoValue = errors.New("no value")

var ErrUnsupportedType = errors.New("unsupported metric type")

func (Counter) Type() string {
	return StringCounterType
}

func (c Counter) ToInt64() int64 {
	return int64(c)
}

func (Gauge) Type() string {
	return StringGaugeType
}

func (g Gauge) ToFloat64() float64 {
	return float64(g)
}

func (c Counter) String() string {
	return strconv.Itoa(int(c))
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}
