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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	newsaggregatorv1 "teamdev.com/go-gator-operator/api/v1"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var k8sClient client.Client

const (
	// serverUri is the link to our news-aggregator
	serverUri = "https://go-gator-svc.go-gator.svc.cluster.local:443/admin/sources"

	// defaultSourceFormat identifies default data format which should be used for new feed
	defaultSourceFormat = "xml"

	// feedFinalizerName is a title of finalizer which will be added to feed object
	// for proper deletion of feed in news aggregator
	feedFinalizerName = "feed.finalizers"
)

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var res ctrl.Result
	var feed newsaggregatorv1.Feed
	var err error

	logger := log.FromContext(ctx)

	err = r.Get(ctx, req.NamespacedName, &feed)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, err
	}

	if feed.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&feed, feedFinalizerName) {
			controllerutil.AddFinalizer(&feed, feedFinalizerName)
			logger.Info("Add Finalizer", feed.Name, feedFinalizerName)

			err = r.Client.Update(ctx, &feed)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(&feed, feedFinalizerName) {
			if _, err = r.handleDelete(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&feed, feedFinalizerName)
			logger.Info("Remove Finalizer", feed.Name, feedFinalizerName)

			err = r.Client.Update(ctx, &feed)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// question: on the review Vitalii said that I should check if key exists in the map
	// but, the logic behind this was that if the key doesn't exist, it means that the feed is new
	// and I should create it. Otherwise, I should update it
	isNew := feed.Status.Conditions[newsaggregatorv1.TypeFeedCreated] == newsaggregatorv1.FeedConditions{} &&
		feed.Status.Conditions[newsaggregatorv1.TypeFeedCreated].Status == false

	if isNew {
		res, err = r.handleCreate(ctx, &feed)
	} else {
		res, err = r.handleUpdate(ctx, &feed)
	}

	if err != nil {
		feed.SetFailedCondition(err.Error(), err.Error())
		return ctrl.Result{}, err
	}

	if isNew {
		feed.SetCreatedCondition("Feed was successfully created")
	} else {
		feed.SetUpdatedCondition("Feed was successfully updated")
	}

	err = r.Client.Status().Update(ctx, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}

	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.Feed{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

// sourceBody is a struct which will be used to generate body for request to the news aggregator.
//
// It contains few key fields which suits perfectly to create/update/delete sources in news aggregator.
type sourceBody struct {
	// Name field describes this source name
	Name string `json:"name"`

	// Format identifies parser which should be used for this particular source
	Format string `json:"format"`

	// Endpoint is a link to this feeds endpoint, where we will go to parse articles.
	Endpoint string `json:"endpoint"`
}

// handleCreate makes a request to the news-aggregator service to create a new feed when a new Feed object is instantiated.
// It constructs a Feed object from the Feed specifications, marshals it to JSON, and sends a POST request with the JSON payload.
// The function handles potential errors in JSON marshalling, request creation, and the HTTP request itself.
// If the server responds with a status other than 201 Created, it attempts to decode and print the server's error message.
func (r *FeedReconciler) handleCreate(ctx context.Context, feed *newsaggregatorv1.Feed) (ctrl.Result, error) {
	source := sourceBody{
		Name:     feed.Spec.Name,
		Format:   defaultSourceFormat,
		Endpoint: feed.Spec.Link,
	}

	sourceData, err := json.Marshal(source)
	if err != nil {
		return ctrl.Result{}, err
	}

	requestBody := bytes.NewBuffer(sourceData)

	req, err := http.NewRequest(http.MethodPost, serverUri, requestBody)
	if err != nil {
		return ctrl.Result{}, err
	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}

	res, err := customClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if res.StatusCode != http.StatusCreated {
		serverError := &serverErr{}

		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, serverError
	}

	err = res.Body.Close()
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// handleUpdate makes a request to the news-aggregator service to update an existing feed when the Feed object is modified.
// It constructs a Feed object from the Feed specifications, marshals it to JSON, and sends a PUT request with the JSON payload.
// This function handles potential errors in JSON marshalling, request creation, and the HTTP request itself.
// If the server responds with a status other than 200 OK, it attempts to decode and print the server's error message.
func (r *FeedReconciler) handleUpdate(ctx context.Context, feed *newsaggregatorv1.Feed) (ctrl.Result, error) {
	source := sourceBody{
		Name:     feed.Spec.Name,
		Format:   defaultSourceFormat,
		Endpoint: feed.Spec.Link,
	}

	sourceData, err := json.Marshal(source)
	if err != nil {
		return ctrl.Result{}, err
	}

	requestBody := bytes.NewBuffer(sourceData)

	req, err := http.NewRequest(http.MethodPut, serverUri, requestBody)
	if err != nil {
		return ctrl.Result{}, err
	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}

	res, err := customClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if res.StatusCode != http.StatusOK {
		serverError := &serverErr{}
		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, serverError
	}

	err = res.Body.Close()
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// handleDelete makes a request to the news-aggregator service to delete an existing feed based on the Feed object.
// It constructs a Feed object from the Feed specifications, marshals it to JSON, and sends a DELETE request with the JSON payload.
// This function handles potential errors in JSON marshalling, request creation, and the HTTP request itself.
// If the server responds with a status other than 200 OK, it attempts to decode and print the server's error message.
func (r *FeedReconciler) handleDelete(ctx context.Context, feed *newsaggregatorv1.Feed) (ctrl.Result, error) {
	source := sourceBody{
		Name:     feed.Spec.Name,
		Format:   defaultSourceFormat,
		Endpoint: feed.Spec.Link,
	}

	if feed.Name == "" {
		return ctrl.Result{}, nil
	}

	sourceData, err := json.Marshal(source)
	if err != nil {
		return ctrl.Result{}, err
	}

	requestBody := bytes.NewBuffer(sourceData)

	req, err := http.NewRequest(http.MethodDelete, serverUri, requestBody)
	if err != nil {
		return ctrl.Result{}, err
	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}

	res, err := customClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if res.StatusCode != http.StatusOK {
		var serverError *serverErr
		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, serverError
	}

	err = res.Body.Close()
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
