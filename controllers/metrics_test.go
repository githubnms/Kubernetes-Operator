package controllers

import (
	"testing"

	autoscalingv1alpha1 "github.com/githubnms/k8s-custom-autoscaler/api/v1alpha1"
)

func TestDecideReplicas(t *testing.T) {
	baseSpec := autoscalingv1alpha1.ScalingPolicySpec{
		ScaleUpThreshold:   100,
		ScaleDownThreshold: 20,
		MinReplicas:        1,
		MaxReplicas:        5,
	}

	tests := []struct {
		name    string
		current int32
		metric  int32
		spec    autoscalingv1alpha1.ScalingPolicySpec
		want    int32
	}{
		{
			name:    "metric above threshold scales up by one",
			current: 2,
			metric:  150,
			spec:    baseSpec,
			want:    3,
		},
		{
			name:    "metric below threshold scales down by one",
			current: 3,
			metric:  5,
			spec:    baseSpec,
			want:    2,
		},
		{
			name:    "metric within normal range makes no change",
			current: 2,
			metric:  50,
			spec:    baseSpec,
			want:    2,
		},
		{
			name:    "never scales above MaxReplicas even if metric is very high",
			current: 5,
			metric:  9999,
			spec:    baseSpec,
			want:    5,
		},
		{
			name:    "never scales below MinReplicas even if metric is very low",
			current: 1,
			metric:  0,
			spec:    baseSpec,
			want:    1,
		},
		{
			name:    "scaling up respects MaxReplicas boundary exactly",
			current: 4,
			metric:  200,
			spec:    baseSpec,
			want:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decideReplicas(tt.current, tt.metric, tt.spec)
			if got != tt.want {
				t.Errorf("decideReplicas(%d, %d, spec) = %d, want %d",
					tt.current, tt.metric, got, tt.want)
			}
		})
	}
}
