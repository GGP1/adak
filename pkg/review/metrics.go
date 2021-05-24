package review

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	totalReviews prometheus.Gauge
	methodCalls  *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "review"
	return metrics{
		totalReviews: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "reviews_total",
			Help:      "Total number of reviews",
		}),
		methodCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "method_calls_total",
			Help:      "Total number of calls per method",
		}, []string{"method"}),
	}
}

func (m metrics) incMethodCalls(method string) {
	m.methodCalls.With(prometheus.Labels{"method": method}).Inc()
}
