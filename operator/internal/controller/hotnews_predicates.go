package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

type FeedStatusConditionPredicate struct {
	predicate.Funcs
}

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
