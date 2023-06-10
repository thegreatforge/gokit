# prom-metrics package

prom-metrics package is a wrapper over prometheus metrics package making easy to instrument the application.
A http server is required with this package as it will self host the metrics rather than pushing them to statsd or otel-collector.

### Usage

```go
// to add a metric

	metrics.GetMonitor().AddMetric(&metrics.Metric{
		Type:        metrics.Counter,
		Name:        "level_debug",
		Description: "Number of debug logs",
		Labels:      []string{"provider"},
	})


// to increment

metrics.GetMonitor().GetMetric("level_debug").Inc([]string{"pod-123-xyz"})
```