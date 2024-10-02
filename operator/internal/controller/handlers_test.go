package controller

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	"testing"
)

func TestHotNewsHandler_UpdateHotNews(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	tests := []struct {
		name             string
		existingHotNews  []newsaggregatorv1.HotNews
		inputObject      client.Object
		expectedRequests []reconcile.Request
		expectedError    bool
		client           client.Client
	}{
		{
			name: "Successfully return reconcile requests for all HotNews",
			existingHotNews: []newsaggregatorv1.HotNews{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews2",
						Namespace: "default",
					},
				},
			},
			client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(
				&newsaggregatorv1.HotNews{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
				&newsaggregatorv1.HotNews{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews2",
						Namespace: "default",
					},
				},
			).Build(),
			inputObject: &newsaggregatorv1.Feed{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			expectedRequests: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
				{
					NamespacedName: types.NamespacedName{
						Name:      "hotnews2",
						Namespace: "default",
					},
				},
			},
			expectedError: false,
		},
		{
			name:            "No HotNews in the namespace",
			existingHotNews: []newsaggregatorv1.HotNews{},
			client:          fake.NewClientBuilder().WithScheme(scheme).Build(),
			inputObject: &newsaggregatorv1.Feed{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			expectedRequests: nil,
			expectedError:    false,
		},
		{
			name: "Error during listing HotNews (invalid namespace)",
			existingHotNews: []newsaggregatorv1.HotNews{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "invalid-namespace",
					},
				},
			},
			inputObject: &newsaggregatorv1.Feed{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			client: fake.NewClientBuilder().WithScheme(scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("listing error")
					},
				}).Build(),
			expectedRequests: nil,
			expectedError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objs := make([]client.Object, len(test.existingHotNews))
			for i, hotNews := range test.existingHotNews {
				objs[i] = hotNews.DeepCopy()
			}

			handler := HotNewsHandler{
				Client: test.client,
			}

			ctx := context.Background()
			requests := handler.UpdateHotNews(ctx, test.inputObject)

			if test.expectedError {
				assert.Nil(t, requests)
			} else {
				assert.Equal(t, test.expectedRequests, requests)
			}
		})
	}
}

func TestConfigMapHandler_ValidateConfigMapFeeds(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	tests := []struct {
		name             string
		existingFeeds    []newsaggregatorv1.Feed
		existingHotNews  []newsaggregatorv1.HotNews
		configMapData    map[string]string
		expectedRequests []reconcile.Request
		expectedError    bool
		client           client.Client
	}{
		{
			name: "Successfully generate reconcile requests for all HotNews",
			existingFeeds: []newsaggregatorv1.Feed{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "feed1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "feed2",
					},
				},
			},
			existingHotNews: []newsaggregatorv1.HotNews{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews2",
						Namespace: "default",
					},
				},
			},
			configMapData: map[string]string{
				"feeds": "feed1,feed2",
			},
			client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(
				&newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Name: "feed1",
					},
				},
				&newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Name: "feed2",
					},
				},
				&newsaggregatorv1.HotNews{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
				&newsaggregatorv1.HotNews{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews2",
						Namespace: "default",
					},
				},
			).Build(),
			expectedRequests: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "hotnews1",
					},
				},
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "hotnews2",
					},
				},
			},
			expectedError: false,
		},
		{
			name:          "Feeds not found in the cluster",
			existingFeeds: []newsaggregatorv1.Feed{},
			existingHotNews: []newsaggregatorv1.HotNews{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
			},
			configMapData: map[string]string{
				"feeds": "nonexistent-feed",
			},
			client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(
				&newsaggregatorv1.HotNews{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hotnews1",
						Namespace: "default",
					},
				},
			).Build(),
			expectedRequests: nil,
			expectedError:    true,
		},
		{
			name: "Error while listing HotNews",
			existingFeeds: []newsaggregatorv1.Feed{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "feed1",
					},
				},
			},
			configMapData: map[string]string{
				"feeds": "feed1",
			},
			client: fake.NewClientBuilder().WithScheme(scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
						return errors.New("listing error")
					},
				}).Build(),
			expectedRequests: nil,
			expectedError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &ConfigMapHandler{
				Client: test.client,
			}

			ctx := context.Background()
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: test.configMapData,
			}
			requests := handler.validateConfigMapFeeds(ctx, configMap)

			if test.expectedError {
				assert.Nil(t, requests)
			} else {
				assert.Equal(t, test.expectedRequests, requests)
			}
		})
	}
}
