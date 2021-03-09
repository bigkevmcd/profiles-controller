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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
)

// ProfileCatalogSourceReconciler reconciles a ProfileCatalogSource object
type ProfileCatalogSourceReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Profiles ProfilesRepository
}

// +kubebuilder:rbac:groups=profiles.weave.works,resources=profilecatalogsources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=profiles.weave.works,resources=profilecatalogsources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=profiles.weave.works,resources=profilecatalogsources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProfileCatalogSource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ProfileCatalogSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("profilecatalogsource", req.NamespacedName)
	log.Info("reconciling catalog source", "req", req)

	var pcs profilesv1alpha1.ProfileCatalogSource
	if err := r.Get(ctx, req.NamespacedName, &pcs); err != nil {
		// This ignores NotFound errors because retrying is unlikely to fix the
		// problem.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Only process embedded profiles just now.

	for n := range pcs.Spec.EmbeddedProfiles {
		r.Profiles.Add(pcs.Spec.EmbeddedProfiles[n])
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileCatalogSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&profilesv1alpha1.ProfileCatalogSource{}).
		Complete(r)
}
