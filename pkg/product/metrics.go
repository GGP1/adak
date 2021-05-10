package product

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	totalProducts prometheus.Gauge
	methodCalls   *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "product"
	return metrics{
		totalProducts: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "products_total",
			Help:      "Total number of products",
		}),
		methodCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "method_calls_total",
			Help:      "Total number of calls per method",
		}, []string{"method"}),
	}
}
