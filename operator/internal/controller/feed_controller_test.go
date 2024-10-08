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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
	"net/http/httptest"
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
					UID:       "123",
					Name:      "ExistingFeedName",
				},
				Spec: newsaggregatorv1.FeedSpec{
					Name: "ExistingFeedName",
					Link: "https://example.com/feed",
				},
			},
		},
	}

	k8sClient := fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()

	type fields struct {
		Client        client.Client
		Scheme        *runtime.Scheme
		serverAddress string
	}
	type args struct {
		ctx context.Context
		req controllerruntime.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockServer *httptest.Server
		want       controllerruntime.Result
		wantErr    bool
		setup      func()
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
			setup: func() {

			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"message": "Feed created successfully"}`))
			})),
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "Performing DELETE of the object",
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
			setup: func() {
				_ = k8sClient.Delete(context.TODO(), &newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "ExistingFeedName",
					},
				})
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message": "Feed created successfully"}`))
			})),
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "Performing UPDATE of the object",
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
			setup: func() {
				_ = k8sClient.Update(context.TODO(), &newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "ExistingFeedName",
					},
					Spec: newsaggregatorv1.FeedSpec{
						Name: "NewFeedName",
						Link: "https://example.com/feed",
					},
				})
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message": "Feed was created successfully"}`))
			})),
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "Error during CREATE of the object",
			fields: fields{
				Client: k8sClient,
				Scheme: nil,
			},
			setup: func() {
				k8sClient.Create(context.TODO(), &newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "ExistingFeedName",
					},
					Spec: newsaggregatorv1.FeedSpec{
						Name: "ExistingFeedName",
						Link: "https://example.com/feed",
					},
				})
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
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "Feed was not created successfully"}`))
			})),
			want:    controllerruntime.Result{},
			wantErr: true,
		},
		{
			name: "Object not found (returning no error)",
			fields: fields{
				Client: k8sClient,
				Scheme: nil,
			},
			args: args{
				ctx: context.TODO(),
				req: controllerruntime.Request{
					NamespacedName: client.ObjectKey{
						Name:      "NonExistentFeed",
						Namespace: "default",
					},
				},
			},
			mockServer: nil,
			setup:      func() {},
			want:       controllerruntime.Result{},
			wantErr:    false,
		},
		{
			name: "Error updating finalizer",
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
			setup: func() {
				feed := &newsaggregatorv1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "ExistingFeedName",
					},
				}
				k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(feed).Build()
				_ = k8sClient.Update(context.TODO(), feed)
			},
			mockServer: nil,
			want:       controllerruntime.Result{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if tt.mockServer != nil {
				tt.fields.serverAddress = tt.mockServer.URL
				defer tt.mockServer.Close()
			}

			r := &FeedReconciler{
				Client:        k8sClient,
				Scheme:        tt.fields.Scheme,
				serverAddress: tt.fields.serverAddress,
			}

			_, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, client.IgnoreNotFound(err))
			}
		})
	}
}

func TestFeedReconciler_handleCreate(t *testing.T) {
	tests := []struct {
		name        string
		feed        *newsaggregatorv1.Feed
		mockServer  *httptest.Server
		expectedErr bool
	}{
		{
			name: "Successful creation",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"message": "Feed created successfully"}`))
			})),
			expectedErr: false,
		},
		{
			name: "JSON marshalling error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: string([]byte{0xff, 0xfe, 0xfd}),
					Link: "http://example.com",
				},
			},
			expectedErr: true,
		},
		{
			name: "HTTP request creation error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: string([]byte{0x7f}),
				},
			},
			expectedErr: true,
		},
		{
			name: "HTTP request performing error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://invalid-url",
				},
			},
			mockServer:  nil,
			expectedErr: true,
		},
		{
			name: "Failed to decode error response",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"invalid": "Feed was not created successfully"}`))
			})),
			expectedErr: true,
		},
		{
			name: "Failed to decode error response",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"invalid": "Feed was not created successfully"}`))
			})),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			if tt.mockServer != nil {
				r.serverAddress = tt.mockServer.URL
				defer tt.mockServer.Close()
			}

			err := r.handleCreate(tt.feed)

			if tt.expectedErr {
				assert.NotNil(t, err)
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
		setup          func(r *FeedReconciler)
		mockServer     *httptest.Server
		expectedResult ctrl.Result
		expectedErr    bool
	}{
		{
			name: "Successful delete",
			setup: func(r *FeedReconciler) {
				err := r.handleCreate(&newsaggregatorv1.Feed{
					Spec: newsaggregatorv1.FeedSpec{
						Name: "Test Feed",
						Link: "http://example.com",
					},
				})
				assert.NotEqual(t, err, "")
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"message": "Feed deleted successfully"}`))
			})),
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
					Name: "Test Feed",
					Link: string([]byte{0xff, 0xfe, 0xfd}),
				},
			},
			setup: func(r *FeedReconciler) {

			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name:  "Invalid feed name",
			setup: func(r *FeedReconciler) {},
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
			setup: func(r *FeedReconciler) {},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name:  "JSON marshalling error",
			setup: func(r *FeedReconciler) {},
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
			name: "Failed to decode error response",
			setup: func(r *FeedReconciler) {

			},
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"invalid": "Feed was not created successfully"}`))
			})),
			expectedResult: ctrl.Result{},
			expectedErr:    false,
		},
		{
			name: "Server returns error",
			setup: func(r *FeedReconciler) {
				err := r.handleCreate(&newsaggregatorv1.Feed{
					Spec: newsaggregatorv1.FeedSpec{
						Name: "Test Feed",
						Link: "http://example.com",
					},
				})
				assert.NotEqual(t, err, "")
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "Feed was not deleted successfully"}`))
			})),
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
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

			if tt.mockServer != nil {
				r.serverAddress = tt.mockServer.URL
				defer tt.mockServer.Close()
			}

			tt.setup(r)

			err := r.handleDelete(tt.feed)

			if tt.expectedErr {
				assert.NotEqual(t, err.Error(), "")
			} else {
				assert.Equal(t, "", err.Error())
			}
		})
	}
}

func TestFeedReconciler_handleUpdate(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		mockServer     *httptest.Server
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
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message": "Feed updated successfully"}`))
			})),
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
		{
			name: "Server returns error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "test Feed",
					Link: "http://example.com",
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "Feed was not updated successfully"}`))
			})),
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			if tt.mockServer != nil {
				r.serverAddress = tt.mockServer.URL
				defer tt.mockServer.Close()
			}

			err := r.handleUpdate(tt.feed)

			if tt.expectedErr {
				assert.NotEqual(t, err.Error(), "")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeedReconciler_SetupWithManager(t *testing.T) {
	schema := runtime.NewScheme()
	assert.Nil(t, newsaggregatorv1.AddToScheme(schema))
	assert.Nil(t, v1.AddToScheme(schema))

	mgr, err := controllerruntime.NewManager(controllerruntime.GetConfigOrDie(), controllerruntime.Options{
		Scheme: schema,
	})
	assert.Nil(t, err)

	type fields struct {
		serverUrl string
		Client    client.Client
		Scheme    *runtime.Scheme
	}
	type args struct {
		mgr       controllerruntime.Manager
		serverUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful setup",
			fields: fields{
				serverUrl: "",
				Client:    mgr.GetClient(),
				Scheme:    mgr.GetScheme(),
			},
			args: args{
				mgr:       mgr,
				serverUrl: "https://go-gator-svc.go-gator.svc.cluster.local:443/news",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{
				serverAddress: tt.fields.serverUrl,
				Client:        tt.fields.Client,
				Scheme:        tt.fields.Scheme,
			}
			err := r.SetupWithManager(tt.args.mgr, tt.args.serverUrl)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.args.serverUrl, r.serverAddress)
		})
	}
}
