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
	return NewMetricsWithRegistry(prometheus.DefaultRegisterer)
}

func NewMetricsWithRegistry(reg prometheus.Registerer) *Metrics {
	factory := promauto.With(reg)

	return &Metrics{
		RequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of HTTP requests processed by Gateway",
			},
			[]string{"method", "endpoint", "status_code", "target_service"},
		),

		RequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "target_service"},
		),

		RequestsActive: factory.NewGauge(
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
