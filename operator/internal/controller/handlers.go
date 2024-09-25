package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

type ConfigMapHandler struct {
	Client client.Client
}

func (r *ConfigMapHandler) UpdateHotNews(ctx context.Context, obj client.Object) []reconcile.Request {
	if _, exists := obj.GetLabels()[newsaggregatorv1.FeedGroupLabel]; !exists {
		return nil
	}

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
