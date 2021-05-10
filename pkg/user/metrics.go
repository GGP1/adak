package user

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	registeredUsers prometheus.Gauge
	methodCalls     *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "user"
	return metrics{
		registeredUsers: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "registered_users_total",
			Help:      "Total number of users registered",
		}),
		methodCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "method_calls_total",
			Help:      "Total number of calls per method",
		}, []string{"method"}),
	}
}
