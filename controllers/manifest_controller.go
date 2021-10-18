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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"github.com/manifest/api/application/v1alpha1"
	"github.com/manifest/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const CheckTime = 30 * time.Second

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	finalizer       = "finalizers.radondb.com/customresource"
)

// ManifestReconciler reconciles a Manifest object
type ManifestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=application.kubesphere.io,resources=manifests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=application.kubesphere.io,resources=manifests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=application.kubesphere.io,resources=manifests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Manifest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ManifestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	customResource := &v1alpha1.Manifest{}
	if err := r.Get(ctx, req.NamespacedName, customResource); err != nil {
		if errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := validateResourceName(customResource); err != nil {
		return ctrl.Result{}, nil
	}

	if customResource.Status.Status == "" {
		customResource.Status.Status = v1alpha1.Creating
		customResource.Status.LastUpdate = metav1.Now()
		err := r.Status().Update(ctx, customResource)
		return reconcile.Result{}, err
	}

	if customResource.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted
		if !sliceutil.HasString(customResource.ObjectMeta.Finalizers, finalizer) {
			customResource.ObjectMeta.Finalizers = append(customResource.ObjectMeta.Finalizers, finalizer)
			err := r.Update(ctx, customResource)
			return reconcile.Result{}, err
		}
	} else {
		// The object is being deleted
		if sliceutil.HasString(customResource.ObjectMeta.Finalizers, finalizer) {
			err := r.deleteCluster(ctx, customResource)
			if err != nil {
				klog.Errorf("delete database cluster error: %s", client.IgnoreNotFound(err).Error())
			}
			customResource.ObjectMeta.Finalizers = sliceutil.RemoveString(customResource.ObjectMeta.Finalizers, func(item string) bool {
				if item == finalizer {
					return true
				}
				return false
			})
			err = r.Update(ctx, customResource)
			return reconcile.Result{}, err
		}
	}
	if customResource.Status.Status == v1alpha1.Creating {
		if err := r.installCluster(ctx, customResource); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	} else if customResource.Status.Version != customResource.Spec.Version {
		if err := r.patchCluster(ctx, customResource); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// check custom resources status
		return r.checkResourceStatus(ctx, customResource)
	}
	klog.V(1).Info("resource name: ", customResource.Name, ", state: ", customResource.Status.Status)
	return ctrl.Result{}, nil
}

func (r *ManifestReconciler) patchCluster(ctx context.Context, resource *v1alpha1.Manifest) error {
	obj, err := getUnstructuredObj(resource)
	if err != nil {
		return err
	}
	err = r.Client.Update(ctx, obj)
	if err != nil {
		klog.V(1).Info(err.Error())
		return client.IgnoreNotFound(err)
	}

	resource.Status.Version = resource.Spec.Version
	err = r.Client.Status().Update(ctx, resource)
	if err != nil {
		resource.Status.Status = v1alpha1.Failed
		err = r.Status().Update(ctx, resource)
	}

	return nil
}

func (r *ManifestReconciler) deleteCluster(ctx context.Context, resource *v1alpha1.Manifest) error {
	klog.V(1).Infof("do delete cluster: %s, %s, %s", resource.Namespace, resource.Name, resource.Spec.Kind)

	obj, err := getUnstructuredObj(resource)
	if err != nil {
		klog.Errorf("get unstructured object error: %s", err.Error())
		return err
	}

	resourceKind := obj.GetKind()
	if resourceKind == "PostgreSQLCluster" {
		// delete pg cluster
		pgCluster := obj.DeepCopy()
		pgCluster.SetKind("Pgcluster")
		pgCluster.SetAPIVersion("radondb.com/v1")
		err = r.Delete(ctx, pgCluster)
		if err != nil {
			klog.Errorf("delete pgcluster resource error: %s", err.Error())
		}
	} else if resourceKind == "Cluster" {
		// delete secret
		err := r.createOrDeleteMysqlClusterPasswordSecret(ctx, resource, true)
		if err != nil {
			klog.Errorf("delete mysql password secret error: %s", err)
			resource.Status.Status = v1alpha1.Error
			err = r.Status().Update(ctx, resource)
			return err
		}
	}

	err = r.Delete(ctx, obj)
	return client.IgnoreNotFound(err)
}

func (r *ManifestReconciler) installCluster(ctx context.Context, resource *v1alpha1.Manifest) error {
	klog.V(1).Infof("install cluster: %s, %s, %s", resource.Namespace, resource.Name, resource.Spec.Kind)
	obj, err := getUnstructuredObj(resource)
	if err != nil {
		return err
	}

	err = r.Create(ctx, obj)
	if err != nil {
		klog.Errorf("create cluster error: %s", err)
		if errors.IsAlreadyExists(err) {
			resource.Status.Status = v1alpha1.AlreadyExists
		} else {
			resource.Status.Status = v1alpha1.Error
		}
		//err = r.Status().Update(ctx, resource)
		//return err
	}

	// mysql 的话创建secret，保存密码
	var clusterStatus string
	resourceKind := obj.GetKind()
	if resourceKind == "Cluster" {
		err := r.createOrDeleteMysqlClusterPasswordSecret(ctx, resource, false)
		if err != nil {
			klog.Errorf("create secret error: %s", err)
			resource.Status.Status = v1alpha1.Error
			err = r.Status().Update(ctx, resource)
			return err
		}
	}

	time.Sleep(500 * time.Millisecond)

	err = r.Get(ctx, types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}, obj)
	if err != nil {
		klog.Error(err, "get unstructured object error.")
		return client.IgnoreNotFound(err)
	}

	if resourceKind == "PostgreSQLCluster" {
		clusterStatus, _ = r.getPgClusterStatus(ctx, obj)
	} else {
		clusterStatus = getUnstructuredObjStatus(obj)
	}
	addObjCondition(obj, resource)

	resource.Status.Status = clusterStatus
	resource.Status.Version = resource.Spec.Version
	switch resource.Kind {
	case v1alpha1.DBTypeClickHouse:
		resource.Spec.AppName = v1alpha1.ClusterAppTypeClickHouse
	case v1alpha1.DBTypePostgreSQL:
		resource.Spec.AppName = v1alpha1.ClusterAPPTypePostgreSQL
	case v1alpha1.DBTypeMysql:
		resource.Spec.AppName = v1alpha1.ClusterAPPTypeMySQL
	default:
		resource.Spec.AppName = ""
	}
	err = r.Client.Status().Update(ctx, resource)
	if err != nil {
		resource.Status.Status = v1alpha1.Failed
		err = r.Status().Update(ctx, resource)
	}
	return nil
}

func (r *ManifestReconciler) createOrDeleteMysqlClusterPasswordSecret(ctx context.Context, resource *v1alpha1.Manifest, delete bool) (err error) {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name + v1alpha1.SuffixSecretName,
			Namespace: resource.Namespace,
		},
	}
	if delete {
		err = r.Client.Delete(ctx, secret)
	} else {
		err = r.Client.Create(ctx, secret)
	}
	return
}

func (r *ManifestReconciler) checkResourceStatus(ctx context.Context, resource *v1alpha1.Manifest) (ctrl.Result, error) {
	klog.V(1).Infof("do check status: %s, %s, %s", resource.Namespace, resource.Name, resource.Spec.Kind)
	obj, err := getUnstructuredObj(resource)
	if err != nil {
		return ctrl.Result{RequeueAfter: CheckTime}, err
	}

	err = r.Get(ctx, types.NamespacedName{
		Namespace: resource.Namespace,
		Name:      resource.Name}, obj)
	if err != nil {
		klog.V(1).Info(err.Error())
		return ctrl.Result{RequeueAfter: CheckTime}, err
	}

	var clusterStatus string
	resourceKind := obj.GetKind()
	if resourceKind == "PostgreSQLCluster" {
		clusterStatus, _ = r.getPgClusterStatus(ctx, obj)
	} else {
		clusterStatus = getUnstructuredObjStatus(obj)
	}
	addObjCondition(obj, resource)
	resource.Status.Status = clusterStatus
	err = r.Client.Status().Update(ctx, resource)
	if err != nil {
		resource.Status.Status = v1alpha1.Failed
		err = r.Status().Update(ctx, resource)
	}
	return ctrl.Result{RequeueAfter: CheckTime}, err
}

func (r *ManifestReconciler) getPgClusterStatus(ctx context.Context, obj *unstructured.Unstructured) (string, error) {
	var pgClusterStatus string
	obj.SetKind("Pgcluster")
	obj.SetAPIVersion("radondb.com/v1")
	err := r.Get(ctx, types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}, obj)
	if err != nil {
		klog.Errorf("get Pgcluster resource error: %s", err.Error())
		return "", client.IgnoreNotFound(err)
	}
	pgClusterStatus = getUnstructuredObjStatus(obj)
	return pgClusterStatus, nil
}
func getUnstructuredObj(resource *v1alpha1.Manifest) (obj *unstructured.Unstructured, err error) {
	obj = &unstructured.Unstructured{}
	_, _, err = decUnstructured.Decode([]byte(resource.Spec.CustomResource), nil, obj)
	if err != nil {
		klog.Errorf("decode unstructured object error: %s", err.Error())
	}
	obj.SetName(resource.Name)
	obj.SetNamespace(resource.Namespace)
	return
}

func getUnstructuredObjStatus(obj *unstructured.Unstructured) string {
	var clusterStatus string
	statusMap, ok := obj.Object["status"].(map[string]interface{})
	if ok {
		clusterStatus, ok = statusMap["state"].(string)
		if ok {
			return clusterStatus
		} else {
			clusterStatus = v1alpha1.ClusterStatusUnknown
		}
	} else {
		clusterStatus = v1alpha1.ClusterStatusUnknown
	}
	return clusterStatus
}

func addObjCondition(obj *unstructured.Unstructured, resource *v1alpha1.Manifest) {
	apiRes, ok := obj.Object["Condition"].(map[string]v1alpha1.ApiResult)
	if ok {
		resource.Status.Condition = apiRes
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManifestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	if r.Scheme == nil {
		r.Scheme = mgr.GetScheme()
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Manifest{}).
		Complete(r)
}

func validateResourceName(resource *v1alpha1.Manifest) error {
	if resource.Spec.Kind == v1alpha1.DBTypeMysql {
		if len(resource.Name) >= 32 {
			return errors.NewBadRequest("The name length can't more than 32 characters")
		}
	} else if resource.Spec.Kind == v1alpha1.DBTypePostgreSQL {
		return errors.NewBadRequest("The name length can't more than 32 characters")
	}
	return nil
}
