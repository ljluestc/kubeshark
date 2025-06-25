package apps

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KubesharkReconciler reconciles Kubeshark deployments
type KubesharkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile handles the reconciliation loop for Kubeshark resources
func (r *KubesharkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Kubeshark deployment", "namespace", req.Namespace, "name", req.Name)

	// Get the pod
	pod := &v1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if this is a Kubeshark pod
	if !isKubesharkPod(pod) {
		return ctrl.Result{}, nil
	}

	// Process Kubeshark pod
	if err := r.processKubesharkPod(ctx, pod); err != nil {
		logger.Error(err, "Failed to process Kubeshark pod", "pod", pod.Name)
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

// isKubesharkPod checks if a pod belongs to Kubeshark
func isKubesharkPod(pod *v1.Pod) bool {
	if pod.Labels == nil {
		return false
	}

	_, isWorker := pod.Labels["app.kubernetes.io/component"]
	_, isHub := pod.Labels["kubeshark/component"]

	return isWorker || isHub
}

// processKubesharkPod handles special requirements for Kubeshark pods
func (r *KubesharkReconciler) processKubesharkPod(ctx context.Context, pod *v1.Pod) error {
	logger := log.FromContext(ctx)

	// Check if the pod needs TLS sniffing capabilities
	if requiresBoringSSL(pod) {
		logger.Info("Kubeshark pod requires BoringSSL support", "pod", pod.Name)

		// Ensure the pod has the necessary security context and volumes
		if err := r.ensureBoringSSLSupport(ctx, pod); err != nil {
			return fmt.Errorf("failed to ensure BoringSSL support: %w", err)
		}
	}

	return nil
}

// requiresBoringSSL checks if a pod requires BoringSSL support
func requiresBoringSSL(pod *v1.Pod) bool {
	// Check if the pod has the BoringSSL annotation
	_, hasAnnotation := pod.Annotations["kubeshark.io/boringssl-enabled"]

	// Check if it's a worker pod (which needs TLS sniffing)
	isWorker := false
	if pod.Labels != nil {
		component, exists := pod.Labels["app.kubernetes.io/component"]
		isWorker = exists && component == "worker"
	}

	return hasAnnotation || isWorker
}

// ensureBoringSSLSupport ensures the pod has the necessary capabilities for BoringSSL
func (r *KubesharkReconciler) ensureBoringSSLSupport(ctx context.Context, pod *v1.Pod) error {
	// Check if we need to update the pod
	needsUpdate := false

	// Check if the pod already has the necessary security context
	if !hasBoringSSLCapabilities(pod) {
		// This would require updating the pod, but pods are immutable
		// We would need to delete and recreate the pod, or use a mutating webhook
		// For this example, we'll just log that an update is needed
		needsUpdate = true
	}

	if needsUpdate {
		// In a real implementation, we might add an annotation to trigger a recreation
		// or use a mutating webhook to modify the pod spec before creation
		log.FromContext(ctx).Info("Pod needs BoringSSL support but cannot be updated directly",
			"pod", pod.Name,
			"namespace", pod.Namespace)
	}

	return nil
}

// hasBoringSSLCapabilities checks if the pod has the necessary capabilities for BoringSSL
func hasBoringSSLCapabilities(pod *v1.Pod) bool {
	// Check if any container has the necessary capabilities
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil &&
			container.SecurityContext.Capabilities != nil {
			caps := container.SecurityContext.Capabilities.Add

			// Check for the necessary capabilities
			hasAllCaps := false
			requiredCaps := []v1.Capability{"SYS_PTRACE", "SYS_ADMIN"}

			for _, reqCap := range requiredCaps {
				found := false
				for _, cap := range caps {
					if cap == reqCap {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}

			if hasAllCaps {
				return true
			}
		}
	}

	return false
}

// SetupWithManager sets up the controller with the Manager
func (r *KubesharkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}
