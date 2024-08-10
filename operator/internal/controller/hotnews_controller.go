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
	"crypto/tls"
	"encoding/json"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

const (
	// serverUrl is a URL to our news aggregator server
	serverUrl = ""

	// pathToKubeConfig is a path to the kubeconfig file
	pathToKubeConfig = ""

	// defaultNamespace is title of the default namespace
	defaultNamespace = ""

	// feedGroupsConfigMapName is a name of the ConfigMap which contains feed groups
	feedGroupsConfigMapName = "feed-group-source"

	// errFeedsAreRequired is an error message which is returned when feeds are not provided
	errFeedsAreRequired = "feeds are required"

	// errKeywordsAreRequired indicates that keywords are required for the request and creation of HotNews object
	errKeywordsAreRequired = "keywords are required"

	// errFailedToConstructRequestUrl is an error message which is returned when failed to construct request URL
	errFailedToConstructRequestUrl = "failed to construct request URL"

	// errFailedToCreateRequest is an error message which is returned when failed to create a new request
	errFailedToCreateRequest = "failed to create a new request"

	// errFailedToSendRequest is an error message which is returned when failed to send a request
	errFailedToSendRequest = "failed to send a request"

	// errFailedToReadResponseBody is an error message which is returned when failed to read response body
	errFailedToReadResponseBody = "failed to read response body"

	// errFailedToUnmarshalResponseBody is an error message which is returned when failed to unmarshal response body
	errFailedToUnmarshalResponseBody = "failed to unmarshal response body"

	// errFailedToGetNews is an error message which is returned when failed to get news
	errFailedToGetNews = "failed to get news"

	// errFailedToCloseResponseBody is an error message which is returned when failed to close response body
	errFailedToCloseResponseBody = "failed to close response body"

	// errFailedToGetConfigMap is an error message which is returned when failed to get ConfigMap
	errFailedToGetConfigMap = "failed to get ConfigMap"

	// errFailedToCreateClientSet is an error message which is returned when failed to create a new client set
	errFailedToCreateClientSet = "failed to create a new client set"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var hotNews newsaggregatorv1.HotNews

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}

	if hotNews.ObjectMeta.CreationTimestamp.IsZero() {
		err = r.handleCreate(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.handleDelete(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = r.handleUpdate(ctx, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

// handleCreate function sends a request to the news aggregator server to retrieve news
// with the specified parameters, and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleCreate(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)
	requestUrl, err := r.constructRequestUrl(
		hotNews.Spec.Keywords,
		hotNews.Spec.DateStart,
		hotNews.Spec.DateEnd,
		hotNews.Spec.Feeds,
	)

	if err != nil {
		logger.Error(err, errFailedToConstructRequestUrl)
		return err
	}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Error(err, errFailedToCreateRequest)
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err, errFailedToSendRequest)
		return err
	}

	if res.StatusCode != http.StatusOK {
		serverErr := &serverError{}
		err = json.NewDecoder(res.Body).Decode(&serverErr)
		if err != nil {
			logger.Error(err, errFailedToUnmarshalResponseBody)
			return err
		}
		return serverErr
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(err, errFailedToCloseResponseBody)
		return err
	}

	return nil
}

// handleUpdate function updates the HotNews object and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleUpdate(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	return nil
}

// handleDelete function deletes the HotNews object and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleDelete(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	return nil
}

// constructRequestUrl function verifies if arguments are correct and constructs a request URL
// to our news aggregator server.
//
// Example:
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&dateEnd=2024-08-06&feeds=abc,bbc
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&feeds=abc,bbc
func (r *HotNewsReconciler) constructRequestUrl(keywords, dateFrom, dateEnd string, feeds []string) (string, error) {
	requestUrl := serverUrl

	if keywords == "" {
		return "", fmt.Errorf(errKeywordsAreRequired)
	}

	if len(feeds) < 1 {
		return "", fmt.Errorf(errFeedsAreRequired)
	}
	requestUrl += fmt.Sprintf("?keywords=%s", keywords)

	var feedStr strings.Builder
	for _, feed := range feeds {
		feedStr.WriteString(feed)
		feedStr.WriteRune(',')
	}

	requestUrl += fmt.Sprintf("&sources=%s", feedStr.String()[:len(feedStr.String())-2])

	if dateFrom != "" {
		requestUrl += fmt.Sprintf("&dateFrom=%s", dateFrom)
	}

	if dateEnd != "" {
		requestUrl += fmt.Sprintf("&dateEnd=%s", dateEnd)
	}

	return requestUrl, nil
}

// getConfigMapData returns all data from config map named feedGroupsConfigMapName in defaultNamespace
func (r *HotNewsReconciler) getConfigMapData(ctx context.Context) ([]string, error) {
	logger := log.FromContext(ctx)

	config := ctrl.GetConfigOrDie()

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err, errFailedToCreateClientSet)
		return nil, err
	}

	configMaps, err := clientSet.CoreV1().ConfigMaps(defaultNamespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logger.Error(err, errFailedToGetConfigMap)
		return nil, err
	}

	var groupData []string

	for _, configMap := range configMaps.Items {
		groupData = append(groupData, configMap.Data[feedGroupsConfigMapName])
	}

	return groupData, nil
}
