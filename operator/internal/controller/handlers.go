package controller

import (
	"context"
	errs "errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

// HotNewsHandler is a struct that is used for triggering HotNews Reconciler,
// whenever a feed which is bound to certain Hot News is being updated.
//
// Fields:
// HotNewsHandler holds a Kubernetes client.
// The client is used to interact with the Kubernetes API, allowing the handler to fetch
// and reconcile HotNews objects when specific events occur.
type HotNewsHandler struct {
	Client client.Client
}

// UpdateHotNews is a method on HotNewsHandler that generates reconcile requests
// when relevant objects are updated in the Kubernetes cluster.
// It retrieves a list of HotNews resources from the namespace of the object
// that triggered the event, and returns a list of reconcile requests for them.
func (r *HotNewsHandler) UpdateHotNews(ctx context.Context, obj client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)

	feed := obj.(*newsaggregatorv1.Feed)

	var configMapList v1.ConfigMapList
	err := r.Client.List(ctx, &configMapList)

	var hotNewsList newsaggregatorv1.HotNewsList
	err = r.Client.List(ctx, &hotNewsList, client.InNamespace(obj.GetNamespace()))
	if err != nil {
		logger.Error(err, "Error during listing hot news:")
		return nil
	}

	var requests []reconcile.Request

	for _, hotNews := range hotNewsList.Items {
		if r.isFeedUsedInHotNews(feed, hotNews, configMapList) {
			requests = append(requests, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: hotNews.Namespace,
					Name:      hotNews.Name,
				},
			})
		}
	}

	return requests
}

// isFeedUsedInHotNews checks whether the provided feed is used in the given HotNews object.
// It returns true if the feed's name is either explicitly listed in HotNews.Spec.Feeds
// or if the feed belongs to a group whose name is listed in HotNews' feed groups,
// which are fetched using the GetFeedGroupNames method with the provided configMapList.
//
// Parameters:
// - feed: The Feed object to check.
// - hotNews: The HotNews object containing feed information.
// - configMapList: A list of ConfigMaps used to resolve feed group names.
//
// Returns:
// - bool: True if the feed is used in the HotNews object, false otherwise.
func (r *HotNewsHandler) isFeedUsedInHotNews(feed *newsaggregatorv1.Feed, hotNews newsaggregatorv1.HotNews,
	configMapList v1.ConfigMapList) bool {
	if slices.Contains(hotNews.Spec.Feeds, feed.Name) {
		return true
	}

	if slices.Contains(hotNews.GetFeedGroupNames(configMapList), feed.Spec.Name) {
		return true
	}

	return false
}

// ConfigMapHandler is a struct responsible for handling and validating ConfigMap objects
// in relation to their association with feeds in a Kubernetes cluster.
//
// Fields:
//   - Client: A Kubernetes client used to interact with the Kubernetes API.
//     It is used to fetch and validate ConfigMap data and its relation to Feed resources.
type ConfigMapHandler struct {
	Client client.Client
}

// validateConfigMapFeeds validates the feeds listed in a ConfigMap's data field, ensuring that
// the feeds exist in the Kubernetes cluster.
// Afterward, it generates reconcile requests for HotNews resources
// associated with the validated feeds.
//
// Parameters:
// - ctx: The context for carrying deadlines, cancellation signals, and other request-scoped values.
// - obj: The Kubernetes ConfigMap object to be validated.
//
// Returns:
// - []reconcile.Request: A list of reconcile requests generated from HotNews resources. Returns nil if validation fails.
func (r *ConfigMapHandler) validateConfigMapFeeds(ctx context.Context, obj client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)

	var configMap *v1.ConfigMap
	var ok bool
	configMap, ok = obj.(*v1.ConfigMap)
	if !ok {
		logger.Error(errs.New("object is not a ConfigMap"), "Invalid object type")
		return nil
	}

	var errList field.ErrorList

	for _, feedGroupValues := range configMap.Data {
		feedNames := strings.Split(feedGroupValues, ",")

		var feed newsaggregatorv1.Feed
		for _, feedName := range feedNames {
			err := r.Client.Get(ctx, client.ObjectKey{
				Name: feedName,
			}, &feed)

			if errors.IsNotFound(err) {
				errList = append(errList, field.Invalid(field.NewPath("data"), feedName, fmt.Sprintf(
					"feed %s is not found. please, create this feed first", feedName)))
			} else {
				logger.Error(err, fmt.Sprintf("Failed to get feed %s from ConfigMap (Client Error)", feedName))
				return nil
			}
		}
	}

	if len(errList) > 0 {
		logger.Error(errList.ToAggregate(), "Feeds were not found in ConfigMap. Please, create them first")
		return nil
	}

	var hotNewsList newsaggregatorv1.HotNewsList

	err := r.Client.List(ctx, &hotNewsList, client.InNamespace(configMap.GetNamespace()))
	if err != nil {
		logger.Error(err, "Error during listing hot news:")
		return nil
	}

	var requests []reconcile.Request

	for _, hotNews := range hotNewsList.Items {
		if r.feedGroupsExistsInMap(hotNews, configMap) {
			requests = append(requests, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: hotNews.Namespace,
					Name:      hotNews.Name,
				},
			})
		}
	}

	return requests
}

// feedGroupsExistsInMap returns true if hotNews contains feedGroups which exists in configMap
func (r *ConfigMapHandler) feedGroupsExistsInMap(hotNews newsaggregatorv1.HotNews, configMap *v1.ConfigMap) bool {
	for _, feedGroupName := range hotNews.Spec.FeedGroups {
		if _, exists := configMap.Data[feedGroupName]; exists {
			return true
		}
	}

	return false
}
