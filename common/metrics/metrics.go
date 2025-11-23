package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	ApiRequests prometheus.Counter
	PRCount     *prometheus.GaugeVec
}

func Init() *Metrics {
	var metrics Metrics

	metrics.ApiRequests = promauto.NewCounter(prometheus.CounterOpts{Name: "api_requests"})
	metrics.PRCount = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "pull_request_count"}, []string{"status"})

	return &metrics
}
