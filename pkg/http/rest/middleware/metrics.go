package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsHandler implements a scrapper middleware.
type MetricsHandler struct {
	requestInFlight *prometheus.GaugeVec
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
}

// NewMetrics initializes the metrics and returns the handler used to scrap.
func NewMetrics() MetricsHandler {
	const ns, sub = "adak", "http"

	basicLabels := []string{"path"}
	httpLabels := []string{"path", "method", "code"}
	sizeBuckets := prometheus.ExponentialBuckets(256, 4, 8)

	return MetricsHandler{
		requestInFlight: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "requests_in_flight",
			Help:      "Number of requests currently handled by this server.",
		}, basicLabels),
		requestCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "requests_total",
			Help:      "Counter of HTTP(S) requests made.",
		}, basicLabels),
		requestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "request_duration_seconds",
			Help:      "Histogram of round-trip request durations.",
			Buckets:   prometheus.DefBuckets,
		}, httpLabels),
		requestSize: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "request_size_bytes",
			Help:      "Total size of the request. Includes body",
			Buckets:   sizeBuckets,
		}, httpLabels),
		responseSize: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "response_size_bytes",
			Help:      "Size of the returned response.",
			Buckets:   sizeBuckets,
		}, httpLabels),
	}
}

// Scrap registers endpoint behavior metrics.
func (m MetricsHandler) Scrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := cleanPath(r.URL.Path)
		labels := prometheus.Labels{"path": path}
		// code is set after the request is intercepted
		httpLabels := prometheus.Labels{"path": path, "method": r.Method, "code": ""}

		m.requestCount.With(labels).Inc()
		inFlight := m.requestInFlight.With(labels)
		inFlight.Inc()

		interceptor := newInterceptor(w)
		next.ServeHTTP(interceptor, r)

		httpLabels["code"] = sanitizeCode(interceptor.statusCode)
		m.requestDuration.With(httpLabels).Observe(time.Since(start).Seconds())
		m.requestSize.With(httpLabels).Observe(float64(computeApproximateRequestSize(r)))
		m.responseSize.With(httpLabels).Observe(float64(interceptor.size))
		inFlight.Dec()
	})
}

// interceptor helps us catch the response status code and response size.
type interceptor struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newInterceptor(w http.ResponseWriter) *interceptor {
	return &interceptor{ResponseWriter: w}
}

// WriteHeader intercepts write header input (status code) and store it in our
// interceptor struct to use it later.
func (h *interceptor) WriteHeader(code int) {
	h.statusCode = code
	h.ResponseWriter.WriteHeader(code)
}

// Write execute the underlying response writer Write and registers the number of bytes written.
func (h *interceptor) Write(b []byte) (int, error) {
	n, err := h.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}
	h.size = n
	return n, nil
}

// https://github.com/prometheus/client_golang/blob/6007b2b5cae01203111de55f753e76d8dac1f529/prometheus/promhttp/instrument_server.go#L298
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

// cleanPath removes the last element from a url path if it's not the only one.
func cleanPath(path string) string {
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	lastIdx := strings.LastIndex(path, "/")
	if lastIdx < 1 {
		return path
	}
	return path[:lastIdx]
}

func sanitizeCode(code int) string {
	switch code {
	case 100:
		return "100"
	case 101:
		return "101"
	case 102:
		return "102"
	case 103:
		return "103"

	case 200, 0:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"
	case 207:
		return "207"
	case 208:
		return "208"
	case 226:
		return "226"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"
	case 308:
		return "308"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"
	case 419:
		return "419"
	case 420:
		return "420"
	case 421:
		return "421"
	case 422:
		return "422"
	case 423:
		return "423"
	case 424:
		return "424"
	case 425:
		return "425"
	case 426:
		return "426"
	case 427:
		return "427"
	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 451:
		return "451"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"
	case 506:
		return "506"
	case 507:
		return "507"
	case 508:
		return "508"
	case 509:
		return "509"
	case 510:
		return "510"
	case 511:
		return "511"

	default:
		return strconv.Itoa(code)
	}
}
