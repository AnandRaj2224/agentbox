package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ExecutionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "agentbox_executions_total",
	Help: "Total number of code executions partitioned by runtime and status.",
}, []string{"runtime", "status"})

var ExecutionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "agentbox_execution_duration_seconds",
	Help:    "Execution duration in seconds partitioned by runtime.",
	Buckets: prometheus.DefBuckets,
}, []string{"runtime"})

func RecordExecution(runtime string, status string) {
	ExecutionsTotal.WithLabelValues(runtime, status).Inc()
}

func RecordDuration(runtime string, durationSeconds float64) {
	ExecutionDuration.WithLabelValues(runtime).Observe(durationSeconds)
}
