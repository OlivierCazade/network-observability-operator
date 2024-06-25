package networkpolicy

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	flowslatest "github.com/netobserv/network-observability-operator/apis/flowcollector/v1beta2"
	"github.com/netobserv/network-observability-operator/controllers/constants"
	"github.com/netobserv/network-observability-operator/controllers/reconcilers"
	"github.com/netobserv/network-observability-operator/pkg/helper"
	"github.com/netobserv/network-observability-operator/pkg/manager"
	"github.com/netobserv/network-observability-operator/pkg/manager/status"
)

type Reconciler struct {
	client.Client
	mgr              *manager.Manager
	status           status.Instance
}

func Start(ctx context.Context, mgr *manager.Manager) error {
	log := log.FromContext(ctx)
	log.Info("Starting Network Policy controller")
	r := Reconciler{
		Client: mgr.Client,
		mgr:    mgr,
		status: mgr.Status.ForComponent(status.NetworkPolicy),
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&flowslatest.FlowCollector{}, reconcilers.IgnoreStatusChange).
		Named("networkPolicy").
		Owns(&networkingv1.NetworkPolicy{}).
		Complete(&r)
}

// Reconcile is the controller entry point for reconciling current state with desired state.
// It manages the controller status at a high level. Business logic is delegated into `reconcile`.
func (r *Reconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	l := log.Log.WithName("networkpolicy") // clear context (too noisy)
	ctx = log.IntoContext(ctx, l)

	// Get flowcollector & create dedicated client
	clh, desired, err := helper.NewFlowCollectorClientHelper(ctx, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get FlowCollector: %w", err)
	} else if desired == nil {
		// Delete case
		return ctrl.Result{}, nil
	}

	r.status.SetUnknown()
	defer r.status.Commit(ctx, r.Client)

	err = r.reconcile(ctx, clh, desired)
	if err != nil {
		l.Error(err, "Network policy reconcile failure")
		// Set status failure unless it was already set
		if !r.status.HasFailure() {
			r.status.SetFailure("NetworkPolicyError", err.Error())
		}
		return ctrl.Result{}, err
	}

	r.status.SetReady()
	return ctrl.Result{}, nil
}

func (r *Reconciler) reconcile(ctx context.Context, clh *helper.Client, desired *flowslatest.FlowCollector) error {
	// log := log.FromContext(ctx)
	ns := helper.GetNamespace(&desired.Spec)
	privilegedNs := ns + constants.EBPFPrivilegedNSSuffix
	authorizedNs := append([]string{privilegedNs}, desired.Spec.NetworkPolicy.AdditionalNamespaces...)
	if desired.Spec.Loki.Mode == flowslatest.LokiModeLokiStack && desired.Spec.Loki.LokiStack.Namespace != "" {
		authorizedNs = append(authorizedNs, desired.Spec.Loki.LokiStack.Namespace)
	}
	npName, desiredNp := buildNetworkPolicy(ns, desired, authorizedNs)
	if err := reconcilers.ReconcileNetworkPolicy(ctx, clh, npName, desiredNp); err != nil {
		return err
	}
	privilegedNpName, desiredPrivilegedNp := buildNetworkPolicy(privilegedNs, desired, append([]string{ns}, desired.Spec.NetworkPolicy.AdditionalNamespaces...))

	err := reconcilers.ReconcileNetworkPolicy(ctx, clh, privilegedNpName, desiredPrivilegedNp)

	return err
}
