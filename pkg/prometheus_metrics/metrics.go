package prometheusmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CacheHitRate A summary for cache hit rates
	CacheHitRate *prometheus.HistogramVec

	// Invocations A counter for invocations
	Invocations *prometheus.CounterVec

	// // Targets A counter for targets
	// Targets *prometheus.CounterVec

	// TargetDurations A gauge for target durations
	//	TargetDurations *prometheus.GaugeVec

	// TargetDurations A summary for target durations
	TargetDurations *prometheus.GaugeVec

	// Tests A counter for tests
	// Tests *prometheus.CounterVec

	// TestDurations A gauge for test durations
	// TestDurations *prometheus.GaugeVec

	// TestDurations A summary for test durations
	TestDurations *prometheus.GaugeVec
)

// RegisterMetrics registers the metrics with Prometheus
func RegisterMetrics() {
	CacheHitRate = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "remote_cache_hit_rate",
		Help:      "Cache hit rate for action cache or content addressable storage across all invocations",
		Buckets:   []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
	}, []string{"CacheType"})

	Invocations = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "invocation_counts",
		Help:      "Number of Invocations Per Host, User and Step Label",
	}, []string{"Host", "User", "StepLabel"})

	TargetDurations = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "slow_target_durations",
		Help:      "Duration of targets that took greater than 2 seconds to complete",
	}, []string{"Target"})

	TestDurations = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "uncached_test_durations",
		Help:      "Durations of tests by target and status",
	}, []string{"Target", "Status", "Strategy"})
}
