package prommetrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricGinRequestTotal string
var metricGinRequest string
var metricGinRequestDuration string
var labelGinRequestHeaderService string

var ginMetricPath = "/metrics"

// Use set gin metrics middleware
// metricsPrefix: metrics prefix, will be used as metrics name prefix.
// labelHeaderService: request header to identify remote service via label.
// metricsRoute: metrics route, default is "/metrics".
func (m *Monitor) Use(r gin.IRoutes, metricsPrefix, labelHeaderService, metricsRoute string) error {

	if metricsPrefix == "" {
		return fmt.Errorf("metricsPrefix can not be empty")
	}

	if labelHeaderService == "" {
		return fmt.Errorf("labelHeaderService can not be empty")
	}

	ginMetricPath = metricsRoute

	metricGinRequestTotal = fmt.Sprintf("%s_gin_request_total", metricsPrefix)
	metricGinRequest = fmt.Sprintf("%s_gin_request", metricsPrefix)
	metricGinRequestDuration = fmt.Sprintf("%s_gin_request_duration", metricsPrefix)
	labelGinRequestHeaderService = labelHeaderService

	err := m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricGinRequestTotal,
		Description: "all the server received request num.",
		Labels:      nil,
	})

	if err != nil {
		return err
	}

	err = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricGinRequest,
		Description: "all the server received request num. with header label key, path, method and status code.",
		Labels:      []string{"remote_service", "path", "method", "status_code"},
	})

	if err != nil {
		return err
	}

	err = m.AddMetric(&Metric{
		Type:        Histogram,
		Name:        metricGinRequestDuration,
		Description: "all the server received request latency. with header label key, path and method.",
		Labels:      []string{"remote_service", "path", "method"},
		Buckets:     []float64{0.1, 0.3, 0.5, 1.2},
	})

	if err != nil {
		return err
	}

	r.Use(m.monitorInterceptor)
	r.GET(ginMetricPath, func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	return nil
}

// monitorInterceptor as gin monitor middleware.
func (m *Monitor) monitorInterceptor(ctx *gin.Context) {
	if ctx.Request.URL.Path == ginMetricPath {
		ctx.Next()
		return
	}
	startTime := time.Now()

	// execute normal process.
	ctx.Next()

	// after request
	m.ginMetricHandle(ctx, startTime)
}

func (m *Monitor) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer
	latency := time.Since(start)

	// set request total
	_ = m.GetMetric(metricGinRequestTotal).Inc(nil)

	headerValue := r.Header.Get(labelGinRequestHeaderService)

	// set request metric
	_ = m.GetMetric(metricGinRequest).Inc([]string{headerValue, ctx.FullPath(), r.Method, strconv.Itoa(w.Status())})

	// set request latency
	_ = m.GetMetric(metricGinRequestDuration).Observe([]string{headerValue, ctx.FullPath(), r.Method}, latency.Seconds())
}
