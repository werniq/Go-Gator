package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

type HotNewsHandler struct {
	Client client.Client
}

func (r *HotNewsHandler) UpdateHotNews(ctx context.Context, obj client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)
	var hotNewsList newsaggregatorv1.HotNewsList

	err := r.Client.List(ctx, &hotNewsList, client.InNamespace(obj.GetNamespace()))
	if err != nil {
		logger.Error(err, "Error during listing hot news:")
		return nil
	}

	var requests []reconcile.Request
	for _, hotNews := range hotNewsList.Items {
		requests = append(requests, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: hotNews.Namespace,
				Name:      hotNews.Name,
			},
		})
	}

	return requests
}
