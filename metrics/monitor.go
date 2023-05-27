package metrics

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Monitor is an object that uses to set server monitor.
type Monitor struct {
	metricPath string
	metrics    map[string]*Metric
}

var monitor *Monitor

// GetMonitor used to get global Monitor object,
// this function returns a singleton object.
func GetMonitor() *Monitor {
	if monitor == nil {
		monitor = &Monitor{
			metricPath: "/metrics",
			metrics:    make(map[string]*Metric),
		}
	}
	return monitor
}

// GetMetric used to get metric object by metric_name.
func (m *Monitor) GetMetric(name string) *Metric {
	if metric, ok := m.metrics[name]; ok {
		return metric
	}
	return &Metric{}
}

// SetMetricPath set metricPath property. metricPath is used for Prometheus
// to get gin server monitoring data.
func (m *Monitor) SetMetricPath(path string) {
	m.metricPath = path
}

// AddMetric add metric object to Monitor.
func (m *Monitor) AddMetric(metric *Metric) error {
	if _, ok := m.metrics[metric.Name]; ok {
		return errors.Errorf("metric '%s' is existed", metric.Name)
	}

	if metric.Name == "" {
		return errors.Errorf("metric name cannot be empty.")
	}
	if f, ok := promTypeHandler[metric.Type]; ok {
		if err := f(metric); err == nil {
			prometheus.MustRegister(metric.vec)
			m.metrics[metric.Name] = metric
			return nil
		}
	}
	return errors.Errorf("metric type '%d' not existed.", metric.Type)
}
