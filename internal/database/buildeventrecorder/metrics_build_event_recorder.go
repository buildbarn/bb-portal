package buildeventrecorder

import (
	"context"
	"math/bits"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	saveBatchDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "buildbarn",
			Subsystem: "portal",
			Name:      "build_event_recorder_save_batch_duration_seconds",
			Help:      "Time spent saving a batch of build events in seconds.",
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 17),
		},
		[]string{"size_bucket", "status"},
	)
	saveBatchSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "buildbarn",
			Subsystem: "portal",
			Name:      "build_event_recorder_save_batch_size",
			Help:      "Number of build events saved in a batch.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 17),
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(saveBatchDuration)
	prometheus.MustRegister(saveBatchSize)
}

func batchSizeLabel(n int) string {
	return strconv.Itoa(1 << bits.Len(uint(n-1)))
}

type metricsBuildEventRecorder struct {
	BuildEventRecorder
}

// NewMetricsBuildEventRecorder decorates a BuildEventRecorder with
// prometheus metrics.
func NewMetricsBuildEventRecorder(buildEventRecorder BuildEventRecorder) BuildEventRecorder {
	return &metricsBuildEventRecorder{
		buildEventRecorder,
	}
}

// SaveBatch keeps track of the time spent saving a batch of build
// events and records metrics for them.
func (r *metricsBuildEventRecorder) SaveBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}
	start := time.Now()
	ret := r.BuildEventRecorder.SaveBatch(ctx, batch)
	status := "ok"
	if ret != nil {
		status = "error"
	}
	duration := time.Since(start)
	saveBatchDuration.WithLabelValues(batchSizeLabel(len(batch)), status).Observe(duration.Seconds())
	saveBatchSize.WithLabelValues(status).Observe(float64(len(batch)))
	return ret
}
