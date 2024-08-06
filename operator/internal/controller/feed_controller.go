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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// SourceBody is a struct which will be used to generate body for request to the news aggregator.
//
// It contains few key fields which suits perfectly to create/update/delete sources in news aggregator.
type SourceBody struct {
	// Name field describes this source name
	Name string `json:"name"`

	// Format identifies parser which should be used for this particular source
	Format string `json:"format"`

	// Endpoint is a link to this feeds endpoint, where we will go to parse articles.
	Endpoint string `json:"endpoint"`
}

const (
	// serverUri is the link to our news-aggregator
	serverUri = "https://10.244.0.17:443/admin/sources"

	// defaultSourceFormat identifies default data format which should be used for new feed
	defaultSourceFormat = "xml"

	// TypeCreated identifies that feed was created
	TypeCreated = "Created"

	// TypeUpdated identifies that feed was updated
	TypeUpdated = "Updated"

	// TypeDeleted identifies that feed was deleted
	TypeDeleted = "Deleted"
)

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=feeds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var res ctrl.Result
	var feed newsaggregatorv1.Feed

	err := r.Get(ctx, req.NamespacedName, &feed)
	switch {
	case errors.IsNotFound(err):
		res, err = r.handleDelete(ctx, &feed)
		if err != nil {
			l.Error(err, "Error while handling delete event ")
			return ctrl.Result{}, err
		}
		return res, nil
	case err != nil:
		return ctrl.Result{}, err
	}

	isNew := feed.Status.Conditions[TypeCreated] != nil

	if isNew {
		res, err = r.handleCreate(ctx, &feed)
	} else {
		res, err = r.handleUpdate(ctx, &feed)
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Client.Status().Update(ctx, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.Feed{}).
		Complete(r)
}

// handleCreate makes a request to the news-aggregator service to create a new feed when a new Feed object is instantiated.
// It constructs a Feed object from the Feed specifications, marshals it to JSON, and sends a POST request with the JSON payload.
// The function handles potential errors in JSON marshalling, request creation, and the HTTP request itself.
// If the server responds with a status other than 201 Created, it attempts to decode and print the server's error message.
func (r *FeedReconciler) handleCreate(ctx context.Context, feed *newsaggregatorv1.Feed) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	source := SourceBody{
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

	l.Info("Successfully executed default client")

	if res.StatusCode != http.StatusCreated {
		l.Info(res.Status)
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

	err = r.updateFeedStatus(ctx, feed, "", true, "Created", "Feed has been created successfully")
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
	source := &SourceBody{
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
		var serverError *serverErr
		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// handleDelete makes a request to the news-aggregator service to delete an existing feed based on the Feed object.
// It constructs a Feed object from the Feed specifications, marshals it to JSON, and sends a DELETE request with the JSON payload.
// This function handles potential errors in JSON marshalling, request creation, and the HTTP request itself.
// If the server responds with a status other than 200 OK, it attempts to decode and print the server's error message.
func (r *FeedReconciler) handleDelete(ctx context.Context, feed *newsaggregatorv1.Feed) (ctrl.Result, error) {
	source := &SourceBody{
		Name: feed.Spec.Name,
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
	}

	return ctrl.Result{}, nil
}

// updateFeedStatus updates the status of the newly-created Feed object with the given status, reason, and message, so that
// feed.Status.Conditions wont be nil and next event will be either update or delete.
func (r *FeedReconciler) updateFeedStatus(ctx context.Context,
	feed *newsaggregatorv1.Feed, t string, status bool, reason string, message string) error {
	condition := &newsaggregatorv1.FeedConditions{
		Status:         status,
		Reason:         reason,
		Message:        message,
		LastUpdateTime: metav1.Now().String(),
	}

	feed.Status.Conditions[t] = condition

	err := r.Client.Status().Update(ctx, feed)
	if err != nil {
		return err
	}

	return nil
}
