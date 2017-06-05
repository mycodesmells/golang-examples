package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsResponseTime prometheus.Summary
)

func init() {
	httpRequestsResponseTime = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "http",
		Name:      "response_time_seconds",
		Help:      "Request response times",
	})

	prometheus.MustRegister(httpRequestsResponseTime)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		httpRequestsResponseTime.Observe(float64(time.Since(start).Seconds()))
	})
}
