package middleware

import (
	"bytes"
	"encoding/json"
	"regexp"
	"time"

	"forum-gateway/handler"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method, path and error type.",
		},
		[]string{"method", "path", "error_type"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds, partitioned by method, path and error type.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "error_type"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// MetricsHandler returns a gin.HandlerFunc that serves the /metrics endpoint for Prometheus scraping.
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return gin.WrapH(h)
}

type metricsBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w metricsBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// classifyError maps business error code to error_type label.
// 0                 → "none"
// 1xxxx             → "system_error" (系统级错误)
// 2xxxx             → "user_error"   (用户非法操作)
// 4xx (HTTP)        → "user_error"   (如 404 路由未匹配)
// 5xx (HTTP) / 其他  → "system_error"
func classifyError(code int) string {
	switch {
	case code == 0:
		return "none"
	case code >= 10000 && code < 20000:
		return "system_error"
	case code >= 20000 && code < 30000:
		return "user_error"
	case code >= 400 && code < 500:
		return "user_error"
	default:
		return "system_error"
	}
}

// Metrics is a middleware that records Prometheus metrics for each HTTP request.
// It captures response body to extract the business error code, then records:
//   - http_requests_total (counter)
//   - http_request_duration_seconds (histogram)
//
// Skipped paths: /swagger/*, /sd/*, /ws
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		if !regexp.MustCompile(`^/api/v1($|/)`).MatchString(path) {
			c.Next()
			return
		}

		if path == "/api/v1/metrics" || path == "/api/v1/ws" || path == "/api/v1/webhook" {
			c.Next()
			return
		}

		method := c.Request.Method
		routePath := c.FullPath()
		if routePath == "" {
			routePath = path
		}

		// Wrap response writer to capture body for error code extraction
		mbw := &metricsBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = mbw

		c.Next()

		duration := time.Since(start).Seconds()

		// Parse response body to get business error code
		code := -1
		var response handler.Response
		if err := json.Unmarshal(mbw.body.Bytes(), &response); err == nil {
			code = response.Code
		}

		errorType := classifyError(code)

		httpRequestsTotal.WithLabelValues(method, routePath, errorType).Inc()
		httpRequestDuration.WithLabelValues(method, routePath, errorType).Observe(duration)
	}
}
