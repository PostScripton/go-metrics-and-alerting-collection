package metrics

import (
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
