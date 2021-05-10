package shop

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	registeredShops prometheus.Gauge
	methodCalls     *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "shop"
	return metrics{
		registeredShops: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "registered_shops_total",
			Help:      "Total number of shops registered",
		}),
		methodCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "method_calls_total",
			Help:      "Total number of calls per method",
		}, []string{"method"}),
	}
}
