/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"os/exec"

	"github.com/DataWorkflowServices/dws/utils/updater"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	nnfv1alpha1 "github.com/nearnodeflash/nnf-integration-test-operator/api/v1alpha1"
)

// IntTestHelperReconciler reconciles a IntTestHelper object
type IntTestHelperReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	WatchNamespace string
}

//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttesthelpers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttesthelpers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttesthelpers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IntTestHelper object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *IntTestHelperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	logger := log.FromContext(ctx)
	var elapsed metav1.Duration
	elapsed.Duration = 0

	helper := &nnfv1alpha1.IntTestHelper{}
	if err := r.Get(ctx, req.NamespacedName, helper); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Don't do anything if we already have a status
	if len(helper.Status.Command) > 0 {
		return ctrl.Result{}, nil
	}

	statusUpdater := updater.NewStatusUpdater[*nnfv1alpha1.IntTestHelperStatus](helper)
	defer func() { err = statusUpdater.CloseWithStatusUpdate(ctx, r.Client.Status(), err) }()

	// Run the command
	// TODO: can this run a script?
	cmd := exec.Command("/bin/bash", "-c", helper.Spec.Command)
	start := metav1.NowMicro()
	output, err := cmd.CombinedOutput()
	elapsed.Duration = metav1.NowMicro().Time.Sub(start.Time)

	logger.Info("Execute", "command", cmd, "output", string(output), "error", err, "elapsed", elapsed.Duration)

	// Update the status with the results
	helper.Status.Command = cmd.String()
	helper.Status.Output = string(output)
	if err != nil {
		helper.Status.Error = err.Error()
	}
	helper.Status.ElapsedTime = elapsed

	return ctrl.Result{}, nil
}

func filterByNamespace(namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return object.GetNamespace() == namespace
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *IntTestHelperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nnfv1alpha1.IntTestHelper{}).
		WithEventFilter(filterByNamespace(r.WatchNamespace)).
		Complete(r)
}
