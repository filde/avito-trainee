package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	ApiRequests prometheus.Counter
}

func Init() *Metrics {
	var metrics Metrics

	metrics.ApiRequests = promauto.NewCounter(prometheus.CounterOpts{Name: "api_requests"})

	return &metrics
}
