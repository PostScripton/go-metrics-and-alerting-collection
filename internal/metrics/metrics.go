package metrics

import (
	"errors"
	"fmt"
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

func (c Counter) Type() string {
	return StringCounterType
}

func (g Gauge) Type() string {
	return StringGaugeType
}

func (c Counter) String() string {
	return strconv.Itoa(int(c))
}

func (g Gauge) String() string {
	return fmt.Sprintf("%f", float64(g))
}
