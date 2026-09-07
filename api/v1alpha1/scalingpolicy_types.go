package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalingPolicySpec defines the desired state of a ScalingPolicy.
// This is the "what the user wants" part — filled in by whoever creates
// the object (e.g. `kubectl apply -f my-policy.yaml`).
type ScalingPolicySpec struct {
	// TargetDeploymentName is the name of the Deployment this policy scales.
	// +kubebuilder:validation:Required
	TargetDeploymentName string `json:"targetDeploymentName"`

	// MetricEndpoint is where the controller reads the current metric value
	// from (e.g. an HTTP endpoint returning a plain number, or later, a
	// Prometheus query). Starting simple: a URL returning JSON like
	// {"value": 42}.
	// +kubebuilder:validation:Required
	MetricEndpoint string `json:"metricEndpoint"`

	// ScaleUpThreshold: if the metric value goes above this, scale up.
	// +kubebuilder:validation:Required
	ScaleUpThreshold int32 `json:"scaleUpThreshold"`

	// ScaleDownThreshold: if the metric value goes below this, scale down.
	// +kubebuilder:validation:Required
	ScaleDownThreshold int32 `json:"scaleDownThreshold"`

	// MinReplicas is the lowest the target Deployment is ever scaled to,
	// even if the metric says to go lower. Protects against scaling to zero
	// unintentionally.
	// +kubebuilder:validation:Minimum=0
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas is the highest the target Deployment is ever scaled to,
	// even if the metric says to go higher. Protects against runaway cost.
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int32 `json:"maxReplicas"`

	// CooldownSeconds is the minimum time to wait between scaling actions,
	// to avoid rapidly flapping replica counts up and down.
	// +kubebuilder:default=30
	CooldownSeconds int32 `json:"cooldownSeconds,omitempty"`
}

// ScalingPolicyStatus reflects the observed/actual state — filled in by
// OUR controller, not the user. This is how the controller reports back
// what it actually did, so anyone can inspect it with `kubectl describe`.
type ScalingPolicyStatus struct {
	// CurrentReplicas is what the controller last observed/set.
	CurrentReplicas int32 `json:"currentReplicas,omitempty"`

	// LastObservedMetricValue is the most recent metric value read.
	LastObservedMetricValue int32 `json:"lastObservedMetricValue,omitempty"`

	// LastScaleTime records when the last scaling action happened —
	// used to enforce the cooldown.
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`

	// Conditions follow the standard Kubernetes pattern for reporting
	// state (e.g. "Ready", "Degraded") in a way other tools can understand.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Target",type=string,JSONPath=`.spec.targetDeploymentName`
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.status.currentReplicas`

// ScalingPolicy is the Schema for the scalingpolicies API.
type ScalingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalingPolicySpec   `json:"spec,omitempty"`
	Status ScalingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScalingPolicyList contains a list of ScalingPolicy — required boilerplate
// so Kubernetes knows how to return multiple objects from a LIST call
// (e.g. `kubectl get scalingpolicies`).
type ScalingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalingPolicy `json:"items"`
}
