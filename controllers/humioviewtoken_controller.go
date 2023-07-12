/*
Copyright 2023 Humio https://humio.com

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
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	corev1alpha1 "github.com/humio/humio-operator/api/v1alpha1"
	humiov1alpha1 "github.com/humio/humio-operator/api/v1alpha1"
	"github.com/humio/humio-operator/pkg/helpers"
	"github.com/humio/humio-operator/pkg/humio"
	"github.com/humio/humio-operator/pkg/kubernetes"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// HumioViewTokenReconciler reconciles a HumioViewToken object
type HumioViewTokenReconciler struct {
	client.Client
	BaseLogger  logr.Logger
	Log         logr.Logger
	HumioClient humio.Client
	Namespace   string
}

//+kubebuilder:rbac:groups=core.humio.com,resources=humioviewtokens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.humio.com,resources=humioviewtokens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.humio.com,resources=humioviewtokens/finalizers,verbs=update

func (r *HumioViewTokenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if r.Namespace != "" {
		if r.Namespace != req.Namespace {
			return reconcile.Result{}, nil
		}
	}

	r.Log = r.BaseLogger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name, "Request.Type", helpers.GetTypeName(r), "Reconcile.ID", kubernetes.RandomString())
	r.Log.Info("Reconciling HumioViewToken")

	// Fetch the HumioViewToken instance
	hvt := &humiov1alpha1.HumioViewToken{}
	err := r.Get(ctx, req.NamespacedName, hvt)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	cluster, err := helpers.NewCluster(ctx, r, hvt.Spec.ManagedClusterName, hvt.Spec.ExternalClusterName, hvt.Namespace, helpers.UseCertManager(), true)
	if err != nil || cluster == nil || cluster.Config() == nil {
		r.Log.Error(err, "unable to obtain humio client config")
		err = r.setState(ctx, humiov1alpha1.HumioViewTokenStateConfigError, hvt)
		if err != nil {
			return reconcile.Result{}, r.logErrorAndReturn(err, "unable to set cluster state")
		}
		return reconcile.Result{RequeueAfter: time.Second * 15}, nil
	}

	// Get current view token
	r.Log.Info("get current view token")
	curToken, err := r.HumioClient.GetViewToken(cluster.Config(), req, hvt)
	if err != nil {
		return reconcile.Result{}, r.logErrorAndReturn(err, "could not check if view token exists")
	}

	emptyToken := humio.ViewToken{}
	if emptyToken == *curToken {
		r.Log.Info("view token doesn't exist. Now adding view token")
		// create token
		_, err := r.HumioClient.AddViewToken(cluster.Config(), req, hvt)
		if err != nil {
			return reconcile.Result{}, r.logErrorAndReturn(err, "could not create view token")
		}
		r.Log.Info("created view token")
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{RequeueAfter: time.Second * 15}, nil
}

func (r *HumioViewTokenReconciler) setState(ctx context.Context, state string, hvt *humiov1alpha1.HumioViewToken) error {
	if hvt.Status.State == state {
		return nil
	}
	r.Log.Info(fmt.Sprintf("setting view token state to %s", state))
	hvt.Status.State = state
	return r.Status().Update(ctx, hvt)
}

func (r *HumioViewTokenReconciler) logErrorAndReturn(err error, msg string) error {
	r.Log.Error(err, msg)
	return fmt.Errorf("%s: %w", msg, err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HumioViewTokenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.HumioViewToken{}).
		Complete(r)
}
