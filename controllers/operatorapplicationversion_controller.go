/*
Copyright 2021.

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

package controllers

import (
	"context"
	sliceutil "github.com/kubesphere/controllers/utils"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	applicationv1alpha1 "github.com/kubesphere/api/application/v1alpha1"
)

const (
	operatorAppVersionFinalizer = "operatorapplication.application.kubesphere.io"
)

// OperatorApplicationVersionReconciler reconciles a OperatorApplicationVersion object
type OperatorApplicationVersionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=application.kubesphere.io,resources=operatorapplicationversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=application.kubesphere.io,resources=operatorapplicationversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=application.kubesphere.io,resources=operatorapplicationversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OperatorApplicationVersion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *OperatorApplicationVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	operatorAppVersion := &applicationv1alpha1.OperatorApplicationVersion{}
	err := r.Client.Get(ctx, req.NamespacedName, operatorAppVersion)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if operatorAppVersion.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted
		if !sliceutil.ContainsString(operatorAppVersion.Finalizers, operatorAppVersionFinalizer) {
			operatorAppVersion.Finalizers = append(operatorAppVersion.Finalizers, operatorAppVersionFinalizer)
			err := r.Update(ctx, operatorAppVersion)
			return ctrl.Result{}, err
		}
	} else {
		// The object is being deleted
		if sliceutil.ContainsString(operatorAppVersion.Finalizers, operatorAppVersionFinalizer) {
			operatorAppVersion.Finalizers = sliceutil.RemoveString(operatorAppVersion.Finalizers, func(item string) bool {
				if item == operatorAppVersionFinalizer {
					return true
				} else {
					return false
				}
			})
			err = r.Update(ctx, operatorAppVersion)
			return ctrl.Result{}, err
		}
	}
	// todo
	/**
	1.更新ownerReferences
	2.关联manifest
	4.删除app时删除相关的manifest
	*/

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OperatorApplicationVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&applicationv1alpha1.OperatorApplicationVersion{}).
		Complete(r)
}
