package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	autoscalingv1alpha1 "github.com/githubnms/k8s-custom-autoscaler/api/v1alpha1"
)

// metricResponse is the expected shape of whatever MetricEndpoint returns.
// Kept intentionally simple for v1 — a real system might later support
// a Prometheus query response instead of this custom format.
type metricResponse struct {
	Value int32 `json:"value"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

// readMetric fetches the current metric value from the configured endpoint.
// Isolated as its own function so it can be swapped out (mock in tests,
// real HTTP call in production) without touching Reconcile at all.
func readMetric(endpoint string) (int32, error) {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return 0, fmt.Errorf("calling metric endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("metric endpoint returned status %d", resp.StatusCode)
	}

	var parsed metricResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, fmt.Errorf("decoding metric response: %w", err)
	}

	return parsed.Value, nil
}

// decideReplicas is pure logic: given the current replica count, the
// latest metric value, and the policy's thresholds/bounds, what should
// the replica count become? No Kubernetes API calls happen in here —
// that's what makes this easy to unit test.
func decideReplicas(current int32, metric int32, spec autoscalingv1alpha1.ScalingPolicySpec) int32 {
	desired := current

	switch {
	case metric > spec.ScaleUpThreshold:
		desired = current + 1
	case metric < spec.ScaleDownThreshold:
		desired = current - 1
	}

	// Enforce safety bounds — never go below min or above max,
	// regardless of what the metric says.
	if desired < spec.MinReplicas {
		desired = spec.MinReplicas
	}
	if desired > spec.MaxReplicas {
		desired = spec.MaxReplicas
	}

	return desired
}
