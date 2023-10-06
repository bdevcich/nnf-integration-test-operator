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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dmv1alpha1 "github.com/NearNodeFlash/nnf-dm/api/v1alpha1"
	nnfv1alpha1 "github.com/nearnodeflash/nnf-integration-test-operator/api/v1alpha1"
)

const (
	finalizer = "inttest.cray.hpe.com"
)

// IntTestManagerReconciler reconciles a IntTestManager object
type IntTestManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttestmanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttestmanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nnf.cray.hpe.com,resources=inttestmanagers/finalizers,verbs=update
//+kubebuilder:rbac:groups=dm.cray.hpe.com,resources=datamovementmanagers,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IntTestManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *IntTestManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconcile start")

	m := &nnfv1alpha1.IntTestManager{}
	if err := r.Get(ctx, req.NamespacedName, m); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the DM worker daemonset
	dmm := &dmv1alpha1.DataMovementManager{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nnf-dm-manager-controller-manager",
			Namespace: "nnf-dm-system",
		},
	}
	if err := r.Get(ctx, client.ObjectKeyFromObject(dmm), dmm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !m.GetDeletionTimestamp().IsZero() {
		log.Info("Removing helper to Data Movement Manager Spec")

		// Remove the container from the Data Movement manager
		for i, c := range dmm.Spec.Template.Spec.Containers {
			if c.Name == m.Spec.Template.Spec.Containers[0].Name {
				dmm.Spec.Template.Spec.Containers = append(
					dmm.Spec.Template.Spec.Containers[:i],
					dmm.Spec.Template.Spec.Containers[i+1:]...)
			}
		}
		if err := r.Update(ctx, dmm); err != nil {
			if !apierrors.IsConflict(err) {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		// TODO: is this correct?
		if controllerutil.ContainsFinalizer(m, finalizer) {
			controllerutil.RemoveFinalizer(m, finalizer)

			if err := r.Update(ctx, m); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(m, finalizer) {
		controllerutil.AddFinalizer(m, finalizer)
		if err := r.Update(ctx, m); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	found := false
	for _, c := range dmm.Spec.Template.Spec.Containers {
		if c.Name == m.Spec.Template.Spec.Containers[0].Name {
			log.Info("helper found in Data Movement Manager Spec")
			found = true
			break
		}
	}

	if !found {
		dmm.Spec.Template.Spec.Containers = append(dmm.Spec.Template.Spec.Containers, m.Spec.Template.Spec.Containers[0])

		log.Info("Adding helper to Data Movement Manager Spec", "template spec", dmm.Spec.Template.Spec, "container", m.Spec.Template.Spec.Containers[0])
		if err := r.Update(ctx, dmm); err != nil {
			if !apierrors.IsConflict(err) {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
	}

	log.Info("Reconcile end", "manager", m)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IntTestManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nnfv1alpha1.IntTestManager{}).
		Complete(r)
}
