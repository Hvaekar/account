package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var httpRequestDurHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: "http",
	Name:      "request_duration_seconds",
	Help:      "The latency of the HTTP requests.",
}, []string{"method", "path", "status_code"})

func MustRegisterMetrics(registerer prometheus.Registerer) {
	registerer.MustRegister(httpRequestDurHistogram)
}

func PrometheusMetrics(skipRoutes ...string) gin.HandlerFunc {
	skip := make(map[string]struct{})
	for _, path := range skipRoutes {
		skip[path] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		c.Next()

		if _, ok := skip[path]; ok {
			return
		}

		httpRequestDurHistogram.
			With(prometheus.Labels{
				"method":      c.Request.Method,
				"path":        path,
				"status_code": strconv.Itoa(c.Writer.Status()),
			}).Observe(time.Since(start).Seconds())
	}
}
