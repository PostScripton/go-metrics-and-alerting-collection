package metrics

const (
	StringCounterType = "counter"
	StringGaugeType   = "gauge"
)

type MetricTyper interface {
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
