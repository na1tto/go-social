package observability

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestsInFlight prometheus.Gauge
	HTTPRequestsDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests processed.",
			},
			[]string{"method", "route", "status"},
		),

		HTTPRequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed.",
			},
		),

		HTTPRequestsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		),
	}
}

func (m *Metrics) Handler() http.Handler {

	metricsRegistry := prometheus.NewRegistry()

	metricsRegistry.MustRegister(
		m.HTTPRequestsTotal,
		m.HTTPRequestsInFlight,
		m.HTTPRequestsDuration,
	)

	metricsHandler := promhttp.HandlerFor(
		metricsRegistry,
		promhttp.HandlerOpts{},
	)

	return metricsHandler
}

func (m *Metrics) HTTPRequestStarted() {
	m.HTTPRequestsInFlight.Inc()
}

func (m *Metrics) HTTPRequestFinished(method string, route string, status string, duration time.Duration) {

	m.HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()

	m.HTTPRequestsDuration.WithLabelValues(method, route, status).Observe(duration.Seconds())

	m.HTTPRequestsInFlight.Dec()

}
