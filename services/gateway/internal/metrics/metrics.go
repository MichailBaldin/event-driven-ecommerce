package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestsActive  prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of HTTP requests processed by Gateway",
			},
			[]string{"method", "endpoint", "status_code", "target_service"},
		),

		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets, // [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]
			},
			[]string{"method", "endpoint", "target_service"},
		),

		RequestsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_http_requests_active",
				Help: "Number of HTTP requests currently being processed",
			},
		),
	}
}

func (m *Metrics) RecordRequest(method, endpoint, statusCode, targetService string, duration time.Duration) {
	m.RequestsTotal.WithLabelValues(method, endpoint, statusCode, targetService).Inc()

	m.RequestDuration.WithLabelValues(method, endpoint, targetService).Observe(duration.Seconds())
}

func (m *Metrics) IncActiveRequests() {
	m.RequestsActive.Inc()
}

func (m *Metrics) DecActiveRequests() {
	m.RequestsActive.Dec()
}
