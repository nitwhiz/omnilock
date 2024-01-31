package prom

import "time"

type Sample struct {
	Value     float64
	Timestamp int64
	Labels    map[string]string
}

func NewSample(value float64, opts ...SampleOption) *Sample {
	s := Sample{
		Value:     value,
		Timestamp: time.Now().UnixMilli(),
		Labels:    map[string]string{},
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}
