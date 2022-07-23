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

type Counter int64

type Gauge float64

var ErrNoValue = errors.New("no value")

var ErrUnsupportedType = errors.New("unsupported metric type")

func (Counter) Type() string {
	return StringCounterType
}

func (Gauge) Type() string {
	return StringGaugeType
}

func (c Counter) String() string {
	return strconv.Itoa(int(c))
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}
