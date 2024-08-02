/*
Copyright 2024.

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
	"k8s.io/apimachinery/pkg/api/errors"
	aggregatorv1 "teamdev.com/go-gator/api/aggregator/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=aggregator.teamdev.com,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.teamdev.com,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.teamdev.com,resources=feeds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var res ctrl.Result
	var feed aggregatorv1.Feed

	err := r.Get(ctx, req.NamespacedName, &feed)
	switch {
	case errors.IsNotFound(err):
		res, err = r.handleDelete(ctx, &feed)
		if err != nil {
			return ctrl.Result{}, err
		}
		return res, nil
	case err != nil:
		return ctrl.Result{}, err
	}

	isNew := len(feed.Status.Conditions) == 0

	if isNew {
		res, err = r.handleCreate(ctx, &feed)
	} else {
		res, err = r.handleUpdate(ctx, &feed)
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
		Complete(r)
}

// handleCreate
func (r *FeedReconciler) handleCreate(ctx context.Context, feed *aggregatorv1.Feed) (ctrl.Result, error) {
	// todo: handle create
	return ctrl.Result{}, nil
}

func (r *FeedReconciler) handleUpdate(ctx context.Context, feed *aggregatorv1.Feed) (ctrl.Result, error) {
	// todo: handle update
	return ctrl.Result{}, nil
}

func (r *FeedReconciler) handleDelete(ctx context.Context, feed *aggregatorv1.Feed) (ctrl.Result, error) {
	// todo: handle deletion and remove all articles associated with this feed
	return ctrl.Result{}, nil
}
