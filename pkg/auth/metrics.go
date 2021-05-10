package auth

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	activeSessions prometheus.Gauge
	totalSessions  prometheus.Counter
}

func initMetrics() metrics {
	const ns, sub = "adak", "auth"
	return metrics{
		activeSessions: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "active_sessions_total",
			Help:      "Total number of active sessions",
		}),
		totalSessions: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "sessions_total",
			Help:      "Total number of sessions",
		}),
	}
}
