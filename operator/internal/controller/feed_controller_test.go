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
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
	controllerruntime "sigs.k8s.io/controller-runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	"testing"
)

func TestFeedReconciler_Reconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)

	existingFeedList := &newsaggregatorv1.FeedList{
		Items: []newsaggregatorv1.Feed{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "ExistingFeedName",
				},
				Spec: newsaggregatorv1.FeedSpec{
					Name: "ExistingFeedName",
					Link: "https://example.com/feed",
				},
			},
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()

	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx context.Context
		req controllerruntime.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    controllerruntime.Result
		wantErr bool
	}{
		{
			name: "Successful Reconcile run",
			fields: fields{
				Client: k8sClient,
				Scheme: nil,
			},
			args: args{
				ctx: context.TODO(),
				req: controllerruntime.Request{
					NamespacedName: client.ObjectKey{
						Name:      "ExistingFeedName",
						Namespace: "default",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "Failed Reconcile due to missing Feed",
			fields: fields{
				Client: k8sClient,
				Scheme: nil,
			},
			args: args{
				ctx: context.TODO(),
				req: controllerruntime.Request{
					NamespacedName: client.ObjectKey{
						Name:      "non-existent-feed",
						Namespace: "non-existent-namespace",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: true,
		},
		{
			name: "Failed Reconcile due to Get error",
			fields: fields{
				Client: k8sClient,
				Scheme: nil,
			},
			args: args{
				ctx: context.TODO(),
				req: controllerruntime.Request{
					NamespacedName: client.ObjectKey{
						Name:      "test-feed",
						Namespace: "default",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			_, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeedReconciler_handleCreate(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		mockClient     func() *http.Client
		expectedResult ctrl.Result
		expectedErr    bool
	}{
		{
			name: "Successful creation",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    false,
		},
		{
			name: "JSON marshalling error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: string([]byte{0xff, 0xfe, 0xfd}),
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name: "HTTP request creation error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: string([]byte{0x7f}),
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}
			_, err := r.handleCreate(context.TODO(), tt.feed)

			if tt.expectedErr {
				assert.NotNil(t, err)
				t.Errorf("Got: %v; Expected: %v", err, tt.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeedReconciler_handleDelete(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		setup          func()
		expectedResult ctrl.Result
		expectedErr    bool
	}{
		{
			name: "Successful delete",
			setup: func() {
				r := &FeedReconciler{}
				_, err := r.handleCreate(context.TODO(), &newsaggregatorv1.Feed{
					Spec: newsaggregatorv1.FeedSpec{
						Name: "Test Feed",
						Link: "http://example.com",
					},
				})
				assert.Nil(t, err)
			},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    false,
		},
		{
			name:  "Invalid feed name",
			setup: func() {},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name:  "Invalid feed endpoint",
			setup: func() {},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "",
					Link: "example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name:  "JSON marshalling error",
			setup: func() {},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: string([]byte{0xff, 0xfe, 0xfd}),
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &FeedReconciler{}

			_, err := r.handleDelete(context.TODO(), tt.feed)

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeedReconciler_handleUpdate(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		expectedResult ctrl.Result
		expectedErr    bool
	}{
		{
			name: "Successful update",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    false,
		},
		{
			name: "JSON marshalling error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: string([]byte{0xff, 0xfe, 0xfd}),
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name: "HTTP request creation error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: string([]byte{0x7f}),
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name: "HTTP request execution error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			_, err := r.handleUpdate(context.TODO(), tt.feed)

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
