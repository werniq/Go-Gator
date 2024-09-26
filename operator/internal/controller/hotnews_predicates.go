package controller

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

// FeedStatusConditionPredicate is a predicate that filters UpdateEvents for Feed objects.
// It checks if the feed has specific status conditions (FeedCreated or FeedDeleted) to determine
// if the update event should be processed.
type FeedStatusConditionPredicate struct {
	predicate.Funcs
}

// Update checks the conditions of the Feed object's status.
// It returns true if either the FeedCreated or FeedDeleted status condition is true.
// Returns false if the new object is nil or if the object is not of type Feed.
func (FeedStatusConditionPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectNew == nil {
		return false
	}
	feed, ok := e.ObjectNew.(*newsaggregatorv1.Feed)
	if !ok {
		return false
	}

	if condition, exists := feed.Status.Conditions[newsaggregatorv1.TypeFeedCreated]; exists && condition.Status {
		return true
	}
	if condition, exists := feed.Status.Conditions[newsaggregatorv1.TypeFeedDeleted]; exists && condition.Status {
		return true
	}

	return false
}

// ConfigMapStatusPredicate is a custom predicate that filters UpdateEvents for ConfigMap objects.
// It processes only ConfigMap objects that contain a specific label (FeedGroupLabel).
type ConfigMapStatusPredicate struct {
	predicate.Funcs
}

// Update checks if the ConfigMap contains the necessary FeedGroupLabel.
// It returns true if the label exists, indicating that the object is relevant for further processing.
// Returns false if the new object is nil or if the object is not of type ConfigMap.
func (ConfigMapStatusPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectNew == nil {
		return false
	}

	configMap, ok := e.ObjectNew.(*v1.ConfigMap)
	if !ok {
		return false
	}

	if _, exists := configMap.GetLabels()[newsaggregatorv1.FeedGroupLabel]; !exists {
		return false
	}

	return true
}
