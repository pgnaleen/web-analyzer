package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_analyzer_requests_total",
			Help: "Total number of requests received",
		})

	pageLoadTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "web_analyzer_page_load_seconds",
			Help:    "Time taken to fetch and analyze web pages",
			Buckets: prometheus.DefBuckets,
		})
)

func InitMetrics() {
	prometheus.MustRegister(requestsTotal, pageLoadTime)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func RecordRequest() {
	requestsTotal.Inc()
}

func ObservePageLoad(duration float64) {
	pageLoadTime.Observe(duration)
}
