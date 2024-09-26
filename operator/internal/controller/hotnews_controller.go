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
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

const (
	// errFailedToConstructRequestUrl error message which is returned when failed to construct request URL
	errFailedToConstructRequestUrl = "failed to construct request URL"

	// errFailedToCreateRequest is returned when failed to create a new request
	errFailedToCreateRequest = "failed to create a new request"

	// errFailedToSendRequest indicates error during sending an HTTP request
	errFailedToSendRequest = "failed to send a request"

	// errFailedToDecodeResBody indicates that error occurred when failed to unmarshal response body
	errFailedToDecodeResBody = "failed to decode response body"

	// errFailedToCloseResponseBody is returned when failed to close response body
	errFailedToCloseResponseBody = "failed to close response body"

	// errWrongFeedGroupName is returned when the feed group name is wrong
	errWrongFeedGroupName = "wrong feed group name, please check the feed group name and try again"
)

// HotNewsReconciler reconciles a HotNews object
// Whenever status of HotNews CRD is updated, it sends a request to the news aggregator server
// to retrieve news with the specified parameters.
// It also watches for changes in the ConfigMap with the feed groups, and in the Feed CRD.
//
// Before sending a request to the news aggregator server, it verifies if the arguments are correct:
// - keywords are provided
// - date range is correct
// - feeds or feed groups are provided, and they exists in news aggregator server
// Then, it constructs a request URL and sends a request to the news aggregator server, parses the response
// and updates the HotNews object.
type HotNewsReconciler struct {
	serverUrl string
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews;feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/status;feeds,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/finalizers;feeds,verbs=update
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the Kubernetes reconciliation loop, which ensures that
// the current state of the cluster matches the desired state specified by the HotNews resource.
//
// This function is triggered when a HotNews object is created, updated, or deleted.
// It performs the following tasks:
//  1. Retrieves the HotNews object and corresponding ConfigMap containing feed group information.
//  2. If the HotNews object lacks a finalizer, it adds one to handle cleanup on deletion.
//  3. If the object is marked for deletion, it removes the associated feed reference and updates the object.
//  4. If the object is active, it processes the HotNews by sending a request to the news aggregator server
//     and updating the object based on the response (handled in processHotNews).
//  5. The function also tracks any changes in the ConfigMap and Feed CRD, ensuring they reflect in the HotNews object.
//
// Any errors encountered during retrieval, processing, or updating are logged,
// and appropriate status updates are made to reflect the success or failure of the operation.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var hotNews newsaggregatorv1.HotNews

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	configMapList, err := r.retrieveConfigMap(ctx, hotNews.Namespace)
	if err != nil {
		updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to retrieve config map", err.Error())
		if updateErr != nil {
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(&hotNews, HotNewsFinalizer) {
		controllerutil.AddFinalizer(&hotNews, HotNewsFinalizer)
		err := r.Client.Update(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if !hotNews.DeletionTimestamp.IsZero() {
		err = r.removeFeedReference(ctx, hotNews, configMapList)
		if err != nil {
			updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to remove feed reference for feeds", err.Error())
			if updateErr != nil {
				return ctrl.Result{}, updateErr
			}

			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(&hotNews, HotNewsFinalizer)
		err = r.Client.Update(ctx, &hotNews)
		if err != nil {
			updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to update hotnews: ", err.Error())
			if updateErr != nil {
				return ctrl.Result{}, updateErr
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	err = r.processHotNews(ctx, &hotNews, configMapList)
	if err != nil {
		updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to process hotnews", err.Error())
		if updateErr != nil {
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, err
	}
	logger.Info("HotNews object has been updated")

	err = r.setFeedReference(ctx, hotNews, configMapList)
	if err != nil {
		updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to set feed reference for hotnews", err.Error())
		if updateErr != nil {
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, err
	}

	err = r.Client.Status().Update(ctx, &hotNews)
	if err != nil {
		updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to update hotnews", err.Error())
		if updateErr != nil {
			return ctrl.Result{}, updateErr
		}

		return ctrl.Result{}, err
	}

	err = r.setSuccessfulStatus(ctx, &hotNews, "HotNews successfully updated",
		"Successfully Reconciled Hot News object")
	if err != nil {
		updateErr := r.setFailedStatus(ctx, &hotNews, "Failed to update hotnews", err.Error())
		if updateErr != nil {
			return ctrl.Result{}, updateErr
		}

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager, and initializes the k8s client
// to work with feedGroup Config Map.
// It Watches for any changes in the ConfigMap with the feed groups, and also watches for changes in Feed CRD.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager, serverUrl string) error {
	r.serverUrl = serverUrl

	hotNewsHandler := &HotNewsHandler{Client: mgr.GetClient()}
	configMapHandler := &ConfigMapHandler{Client: mgr.GetClient()}

	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&newsaggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(hotNewsHandler.UpdateHotNews),
			builder.WithPredicates(FeedStatusConditionPredicate{}),
		).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(configMapHandler.UpdateHotNews),
			builder.WithPredicates(ConfigMapStatusPredicate{}),
		).
		Complete(r)
}

// article struct is used to parse and represent individual news articles received from the news aggregator server.
// The struct fields are annotated with both JSON and XML tags to allow parsing of response data in either format.
// Each article contains details such as the title, publication date, description, publisher/source, and a link to the full article.
type article struct {
	// Title The title of the news article.
	Title string `json:"title" xml:"title"`

	// PubDate The date and time when the article was published.
	PubDate string `json:"publishedAt" xml:"pubDate"`

	// Description A brief summary or description of the article's content.
	Description string `json:"description" xml:"description"`

	// Publisher The name of the source or publisher of the article.
	Publisher string `xml:"source" json:"Publisher"`

	// Link The URL link to the full article.
	Link string `json:"url" xml:"link"`
}

// processHotNews updates the given HotNews object by sending an HTTP GET request to a constructed URL
// and processing the response.
//
// The function constructs the request URL based on the HotNews and ConfigMap data,
// sends the request to an external server, and handles the server response.
//
// If the response is successful, it decodes the JSON response body containing articles,
// processes the data (e.g., titles and total count), and updates the HotNews object's status with this information.
//
// Any errors during the process, such as URL construction, request creation, sending, decoding, or closing
// the response body, are logged and returned as errors.
func (r *HotNewsReconciler) processHotNews(ctx context.Context, hotNews *newsaggregatorv1.HotNews, configMapList v1.ConfigMapList) error {
	logger := log.FromContext(ctx)
	logger.Info("handling update")

	requestUrl, err := r.constructRequestUrl(ctx, hotNews, configMapList)
	if err != nil {
		logger.Error(err, errFailedToConstructRequestUrl)
		return err
	}
	logger.Info(requestUrl)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		logger.Error(err, errFailedToCreateRequest)
		return err
	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}

	res, err := customClient.Do(req)
	if err != nil {
		logger.Error(err, errFailedToSendRequest)
		return err
	}

	if res.StatusCode != http.StatusOK {
		serverError := &serverErr{}
		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			logger.Error(err, errFailedToDecodeResBody)
			return err
		}
		return serverError
	}

	var articles struct {
		TotalNews int       `json:"totalAmount"`
		News      []article `json:"news"`
	}

	err = json.NewDecoder(res.Body).Decode(&articles)
	if err != nil {
		logger.Error(err, errFailedToDecodeResBody)
		return err
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(err, errFailedToCloseResponseBody)
		return err
	}

	var articlesTitles []string
	for _, a := range articles.News {
		articlesTitles = append(articlesTitles, a.Title)
	}
	logger.Info("Total amount of news", "totalAmount", articles.TotalNews)

	hotNews.SetStatus(articles.TotalNews, requestUrl, articlesTitles)

	logger.Info("HotNews.processHotNews has been successfully executed")
	logger.Info("HotNews object", "HotNews", hotNews)

	return nil
}

// constructRequestUrl function verifies if arguments are correct and constructs a request URL
// to our news aggregator server.
//
// Example:
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&dateEnd=2024-08-06&sources=abc,bbc
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&sources=abc,bbc
func (r *HotNewsReconciler) constructRequestUrl(ctx context.Context, hotNews *newsaggregatorv1.HotNews,
	configMapList v1.ConfigMapList) (string, error) {
	var requestUrl strings.Builder

	requestUrl.WriteString(r.serverUrl)
	var keywordsStr strings.Builder
	for _, keyword := range hotNews.Spec.Keywords {
		keywordsStr.WriteString(keyword)
		keywordsStr.WriteRune(',')
	}
	requestUrl.WriteString("?keywords=" + keywordsStr.String()[:len(keywordsStr.String())-1])

	var feedStr strings.Builder
	if hotNews.Spec.FeedGroups != nil {
		feedGroupsStr, err := r.processFeedGroups(hotNews, configMapList)
		if err != nil {
			return "", err
		}
		feedStr.WriteString(feedGroupsStr)
	} else {
		feedsStr := r.processFeeds(hotNews.Spec)
		feedStr.WriteString(feedsStr)
	}

	requestUrl.WriteString("&sources=" + feedStr.String())

	if hotNews.Spec.DateStart != "" {
		requestUrl.WriteString("&dateFrom=" + hotNews.Spec.DateStart)
	}

	if hotNews.Spec.DateEnd != "" {
		requestUrl.WriteString("&dateEnd=" + hotNews.Spec.DateEnd)
	}

	return requestUrl.String(), nil
}

// setSuccessfulStatus checks if the condition should be of type "Created" or "Updated".
// It examines the current status of the HotNews object to determine if it has been successfully created before.
// If it has been created, it updates the condition to "Updated"; otherwise, it sets the condition to "Created".
// After setting the appropriate condition, it updates the HotNews object in the Kubernetes cluster.
//
// Parameters:
// - ctx: A context to manage request cancellation and timeouts.
// - hotNews: The HotNews resource whose status needs to be updated.
// - reason: A short string indicating the reason for the status change.
// - message: A descriptive message explaining the status update.
//
// Returns:
// - error: If the update operation on the HotNews object fails, an error is returned. Otherwise, it returns nil.
func (r *HotNewsReconciler) setSuccessfulStatus(ctx context.Context, hotNews *newsaggregatorv1.HotNews,
	reason, message string) error {
	var condition = newsaggregatorv1.TypeHotNewsCreated
	if _, exists := hotNews.Status.Conditions[newsaggregatorv1.HotNewsSuccessfullyCreated]; exists {
		condition = newsaggregatorv1.TypeHotNewsUpdated
	}

	hotNews.SetCondition(condition, newsaggregatorv1.StatusSuccess, reason, message)
	err := r.Client.Update(ctx, hotNews)
	if err != nil {
		return err
	}

	return nil
}

// setFailedStatus sets the status condition of the HotNews resource to "Error".
// It is used when the operation on the HotNews object has failed, and the failure needs to be reflected
// in the status of the resource.
//
// Parameters:
// - ctx: A context to manage request cancellation and timeouts.
// - hotNews: The HotNews resource whose status needs to be updated.
// - reason: A short string indicating the reason for the failure.
// - message: A descriptive message explaining the failure.
//
// Returns:
// - error: If the update operation on the HotNews object fails, an error is returned. Otherwise, it returns nil.
func (r *HotNewsReconciler) setFailedStatus(ctx context.Context, hotNews *newsaggregatorv1.HotNews,
	reason, message string) error {
	hotNews.SetCondition(newsaggregatorv1.HotNewsError, newsaggregatorv1.StatusError, reason, message)
	err := r.Client.Update(ctx, hotNews)
	if err != nil {
		return err
	}

	return nil
}

// setFeedReference sets the owner references for each Feed in the HotNewsSpec.Feeds array.
//
// This function is used to set the feed reference for the feed. It determines if feeds or feedGroups should be used as
// reference, because we support only one of these values.
// It initializes array of feeds which will be passed further for setting owner reference. By default, it is array of
// spec.Feeds, but if feedGroups were not nil (meaning that feeds were nil), it retrieves feed names with
// HotNews method GetFeedGroupNames.
func (r *HotNewsReconciler) setFeedReference(ctx context.Context, hotNews newsaggregatorv1.HotNews,
	configMapList v1.ConfigMapList) error {
	var feeds = hotNews.Spec.Feeds
	if hotNews.Spec.FeedGroups != nil {
		feeds = hotNews.GetFeedGroupNames(configMapList)
	}

	err := r.setOwnerReferenceForFeeds(ctx, hotNews, feeds)
	if err != nil {
		return err
	}

	return nil
}

// setOwnerReferenceForFeeds sets the owner references for each Feed in the given array containing feed names.
//
// This method iterates over the list of feeds specified in the `HotNews` custom resource and ensures
// that each `Feed` resource has the current `HotNews` object set as its owner. The purpose of this
// owner reference is to establish a parent-child relationship between `HotNews` and its associated
// feeds, which enables Kubernetes' garbage collection to automatically delete orphaned feeds
// when the `HotNews` object is deleted.
//
// Parameters:
// - ctx: A context to control cancellation signals and request-scoped values across API boundaries.
// - hotNews: The `HotNews` resource containing metadata and spec about the feeds.
// - feeds: A list of feed names (strings) to update with owner references.
//
// Returns:
//   - error: Returns an error if there's an issue retrieving or updating a feed.
//     Returns an aggregated error if multiple feeds cannot be found.
func (r *HotNewsReconciler) setOwnerReferenceForFeeds(ctx context.Context, hotNews newsaggregatorv1.HotNews, feeds []string) error {
	var errList field.ErrorList

	ownerRef := metav1.NewControllerRef(&hotNews, newsaggregatorv1.GroupVersion.WithKind("HotNews"))

	for _, feedName := range feeds {
		feed := &newsaggregatorv1.Feed{}
		err := r.Client.Get(ctx, client.ObjectKey{
			Name:      feedName,
			Namespace: hotNews.Namespace,
		}, feed)
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				errList = append(errList, field.Invalid(field.NewPath("spec.feeds").Child(feedName), feedName, "feed not found"))
			} else {
				return err
			}
			continue
		}

		if !hasOwnerReference(feed, ownerRef) {
			feed.SetOwnerReferences(append(feed.GetOwnerReferences(), *ownerRef))

			err = r.Update(ctx, feed)
			if err != nil {
				return err
			}
		}
	}

	if len(errList) > 0 {
		return errList.ToAggregate()
	}

	return nil
}

// hasOwnerReference checks if the given Feed already has the provided OwnerReference.
//
// This utility function helps prevent duplication of the owner reference. It compares the UID of the existing
// owner references in the `Feed` object against the `UID` of the provided owner reference.
//
// Parameters:
// - feed: The feed resource to check.
// - ownerRef: The owner reference to check for in the feed's list of owner references.
//
// Returns:
// - bool: Returns `true` if the feed already has the given owner reference, otherwise returns `false`.
func hasOwnerReference(feed *newsaggregatorv1.Feed, ownerRef *metav1.OwnerReference) bool {
	for _, ref := range feed.GetOwnerReferences() {
		if ref.UID == ownerRef.UID {
			return true
		}
	}
	return false
}

// removeFeedReference removes the owner references for each Feed in the HotNewsSpec.Feeds array.
//
// This function is designed to handle the cleanup process when a `HotNews` resource is deleted. It ensures that the
// owner references for the given `HotNews` resource are removed from each associated feed.
//
// Parameters:
//   - ctx: A context to control cancellation signals and request-scoped values.
//   - hotNews: The `HotNews` object containing details about the feeds to update.
//   - configMapList: A list of ConfigMaps that may define feed groups. This is used if `FeedGroups`
//     are defined in the `HotNewsSpec`.
//
// Returns:
// - error: Returns an error if any feed fails to be updated.
func (r *HotNewsReconciler) removeFeedReference(ctx context.Context, hotNews newsaggregatorv1.HotNews,
	configMapList v1.ConfigMapList) error {
	var feeds = hotNews.Spec.Feeds
	if hotNews.Spec.FeedGroups != nil {
		feeds = hotNews.GetFeedGroupNames(configMapList)
	}

	err := r.removeOwnerReferenceFromFeeds(ctx, &hotNews, feeds)
	if err != nil {
		return err
	}

	return nil
}

// removeOwnerReferenceFromFeeds removes the owner references for each Feed in the given feeds array.
//
// This method iterates through the list of feeds provided, removes any existing owner references
// pointing to the given `HotNews` resource, and updates the feed resource in the cluster.
//
// Parameters:
// - ctx: A context to control cancellation signals and request-scoped values.
// - hotNews: A pointer to the `HotNews` resource from which the owner reference is being removed.
// - feeds: A list of feed names whose owner references will be cleared.
//
// Returns:
// - error: Returns an aggregated error if one or more feeds could not be found or updated.
func (r *HotNewsReconciler) removeOwnerReferenceFromFeeds(ctx context.Context, hotNews *newsaggregatorv1.HotNews, feeds []string) error {
	var errList field.ErrorList

	for _, feedName := range feeds {
		feed := &newsaggregatorv1.Feed{}
		err := r.Client.Get(ctx, client.ObjectKey{
			Namespace: hotNews.Namespace,
			Name:      feedName,
		}, feed)
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				errList = append(errList, field.Invalid(field.NewPath("spec.feeds").Child(feedName), feedName, "feed not found"))
			} else {
				return err
			}
			continue
		}

		feed.SetOwnerReferences([]metav1.OwnerReference{})

		err = r.Client.Update(ctx, feed)
		if err != nil {
			return err
		}
	}

	if errList != nil {
		return errList.ToAggregate()
	}

	return nil
}

// processFeeds concatenates the feed sources from the HotNewsSpec.Feeds array into a single
// comma-separated string. It iterates through each feed in the spec and appends it to a
// string builder, followed by a comma. The resulting string is returned without the trailing comma.
func (r *HotNewsReconciler) processFeeds(spec newsaggregatorv1.HotNewsSpec) string {
	var sourcesBuilder strings.Builder

	for _, feed := range spec.Feeds {
		sourcesBuilder.WriteString(feed)
		sourcesBuilder.WriteRune(',')
	}

	return sourcesBuilder.String()[:len(sourcesBuilder.String())-1]
}

// processFeedGroups processes the feed groups defined in the HotNews object by fetching the corresponding
// feeds from the given ConfigMap list. It checks if each feed group exists in the ConfigMap's data,
// and if found, it concatenates the feed sources into a comma-separated string. If a feed group is not found
// in any ConfigMap, it returns an error indicating the wrong feed group name.
func (r *HotNewsReconciler) processFeedGroups(hotNews *newsaggregatorv1.HotNews,
	configMapList v1.ConfigMapList) (string, error) {
	var sourcesBuilder strings.Builder

	for _, feedGroup := range hotNews.Spec.FeedGroups {
		for _, configMap := range configMapList.Items {
			if _, ok := configMap.Data[feedGroup]; !ok {
				return "", fmt.Errorf(errWrongFeedGroupName)
			} else {
				sourcesBuilder.WriteString(configMap.Data[feedGroup])
				sourcesBuilder.WriteRune(',')
			}
		}
	}

	return sourcesBuilder.String()[:len(sourcesBuilder.String())-1], nil
}

// retrieveConfigMap retrieves all config maps from given namespace, based on FeedGroupLabel.
// If config map has this label, it will be retrieved.
//
// This function will be used on the start of HotNewsReconciliation loop, since this config maps will be used often.
// Returns an error either if failed to construct label Requirement, or during listing of config map.
func (r *HotNewsReconciler) retrieveConfigMap(ctx context.Context, namespace string) (v1.ConfigMapList, error) {
	s, err := labels.NewRequirement(newsaggregatorv1.FeedGroupLabel, selection.Exists, nil)
	if err != nil {
		return v1.ConfigMapList{}, err
	}

	var configMaps v1.ConfigMapList
	err = r.Client.List(ctx, &configMaps, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*s),
		Namespace:     namespace,
	})

	if err != nil {
		return v1.ConfigMapList{}, err
	}

	return configMaps, nil
}
