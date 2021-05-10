package account

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	methodCalls *prometheus.CounterVec
}

func initMetrics() metrics {
	const ns, sub = "adak", "account"
	return metrics{
		methodCalls: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "method_calls_total",
			Help:      "Total number of calls per method",
		}, []string{"method"}),
	}
}
