package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// These are custom Prometheus metrics for OUR operator's own behavior —
// separate from the business metric we read to make scaling decisions.
// controller-runtime already exposes a /metrics endpoint by default;
// registering these adds our own operator-specific data to that endpoint.
var (
	currentReplicasGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "scalingpolicy_current_replicas",
			Help: "Current replica count last set by a ScalingPolicy.",
		},
		[]string{"policy", "namespace", "target_deployment"},
	)

	observedMetricGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "scalingpolicy_last_observed_metric",
			Help: "Last metric value observed by a ScalingPolicy.",
		},
		[]string{"policy", "namespace"},
	)

	scaleEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scalingpolicy_scale_events_total",
			Help: "Total number of scaling actions taken, by direction.",
		},
		[]string{"policy", "namespace", "direction"}, // direction: "up" or "down"
	)
)

// init registers our custom metrics with controller-runtime's default
// Prometheus registry, so they're exposed on the same /metrics endpoint
// the manager already serves.
func init() {
	metrics.Registry.MustRegister(currentReplicasGauge, observedMetricGauge, scaleEventsTotal)
}

// recordScalingObservation updates the gauges after every reconcile,
// regardless of whether a scaling action happened.
func recordScalingObservation(policyName, namespace, target string, replicas, metricValue int32) {
	currentReplicasGauge.WithLabelValues(policyName, namespace, target).Set(float64(replicas))
	observedMetricGauge.WithLabelValues(policyName, namespace).Set(float64(metricValue))
}

// recordScaleEvent increments the counter whenever a scaling action
// actually happens — separate from the gauges above, since counters
// answer "how many times has this happened over time," which gauges can't.
func recordScaleEvent(policyName, namespace, direction string) {
	scaleEventsTotal.WithLabelValues(policyName, namespace, direction).Inc()
}
