package controllers

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	autoscalingv1alpha1 "github.com/githubnms/k8s-custom-autoscaler/api/v1alpha1"
)

// ScalingPolicyReconciler reconciles a ScalingPolicy object.
// This struct holds everything the reconcile loop needs to do its job:
// a client to read/write Kubernetes objects, and the type registry (Scheme).
type ScalingPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling.zhero.dev,resources=scalingpolicies,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=autoscaling.zhero.dev,resources=scalingpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

// Reconcile is called every time a ScalingPolicy is created, updated, or
// (indirectly, via requeue) on a timer. This is THE core loop — it always
// follows the same pattern: read current state, compare to desired state,
// take one small corrective action, return.
func (r *ScalingPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Fetch the ScalingPolicy that triggered this reconcile.
	var policy autoscalingv1alpha1.ScalingPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		if apierrors.IsNotFound(err) {
			// Object was deleted — nothing to do, don't requeue.
			return ctrl.Result{}, nil
		}
		// Some other error reading it — requeue and try again.
		return ctrl.Result{}, err
	}

	logger.Info("reconciling ScalingPolicy",
		"name", policy.Name,
		"targetDeployment", policy.Spec.TargetDeploymentName)

	// 2. Fetch the target Deployment this policy is supposed to scale.
	var deploy appsv1.Deployment
	deployKey := types.NamespacedName{
		Namespace: policy.Namespace,
		Name:      policy.Spec.TargetDeploymentName,
	}
	if err := r.Get(ctx, deployKey, &deploy); err != nil {
		// Target doesn't exist yet — not a fatal error, just wait and retry.
		logger.Info("target deployment not found, retrying later", "target", deployKey)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// 3. Enforce cooldown: skip this cycle if we scaled too recently.
	cooldown := time.Duration(policy.Spec.CooldownSeconds) * time.Second
	if policy.Status.LastScaleTime != nil &&
		time.Since(policy.Status.LastScaleTime.Time) < cooldown {
		remaining := cooldown - time.Since(policy.Status.LastScaleTime.Time)
		logger.Info("in cooldown, skipping scaling decision", "remaining", remaining)
		return ctrl.Result{RequeueAfter: remaining}, nil
	}

	// 4. Read the current metric value.
	metricValue, err := readMetric(policy.Spec.MetricEndpoint)
	if err != nil {
		logger.Error(err, "failed to read metric, will retry")
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}

	// 5. Decide the desired replica count based on the metric.
	currentReplicas := int32(1)
	if deploy.Spec.Replicas != nil {
		currentReplicas = *deploy.Spec.Replicas
	}
	desiredReplicas := decideReplicas(currentReplicas, metricValue, policy.Spec)

	// 6. If the decision changed anything, apply it.
	if desiredReplicas != currentReplicas {
		logger.Info("scaling decision made",
			"from", currentReplicas, "to", desiredReplicas, "metric", metricValue)

		deploy.Spec.Replicas = &desiredReplicas
		if err := r.Update(ctx, &deploy); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update deployment replicas: %w", err)
		}

		direction := "up"
		if desiredReplicas < currentReplicas {
			direction = "down"
		}
		recordScaleEvent(policy.Name, policy.Namespace, direction)

		now := metav1.Now()
		policy.Status.LastScaleTime = &now
	}

	// 7. Always refresh status with the latest observed values, and
	// update our Prometheus gauges to match, whether or not a scale
	// action happened this cycle — dashboards need continuous data.
	policy.Status.CurrentReplicas = desiredReplicas
	policy.Status.LastObservedMetricValue = metricValue
	recordScalingObservation(policy.Name, policy.Namespace, policy.Spec.TargetDeploymentName,
		desiredReplicas, metricValue)
	if err := r.Status().Update(ctx, &policy); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update ScalingPolicy status: %w", err)
	}

	// 8. Re-check again after a fixed interval, even if nothing changed —
	// metrics can change without any Kubernetes event happening.
	return ctrl.Result{RequeueAfter: 20 * time.Second}, nil
}

// SetupWithManager wires this reconciler into the controller-runtime
// manager, telling it which resource type to watch.
func (r *ScalingPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1alpha1.ScalingPolicy{}).
		Complete(r)
}
