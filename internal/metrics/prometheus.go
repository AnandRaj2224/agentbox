package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ExecutionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{}, []string{"runtime", "status"})
var ExecutionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Buckets: prometheus.DefBuckets,
}, []string{"runtime"})
