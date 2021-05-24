package ordering

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	totalOrders *prometheus.CounterVec
	methodCalls *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "ordering"
	return metrics{
		totalOrders: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "orders_total",
			Help:      "Total number of orders",
		}, []string{"status"}),
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
