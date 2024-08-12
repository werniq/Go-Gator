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
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
	"reflect"
	controllerruntime "sigs.k8s.io/controller-runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	newsaggregatorv1 "teamdev.com/go-gator-operator/api/v1"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := newsaggregatorv1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
}

func TestFeedReconciler_Reconcile(t *testing.T) {
	feed := &newsaggregatorv1.Feed{}
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
				Client: fake.NewFakeClient(feed),
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
			wantErr: false,
		},
		{
			name: "Failed Reconcile due to missing Feed",
			fields: fields{
				Client: fake.NewFakeClient(),
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
				Client: fake.NewFakeClient(),
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
			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
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
		expectedErr    error
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
			expectedErr:    nil,
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
			expectedErr:    errors.New("json: error calling MarshalJSON"),
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
			expectedErr:    errors.New("http: invalid character"),
		},
		{
			name: "HTTP request execution error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("dial tcp: lookup"),
		},
		{
			name: "Non-201 status code",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("Bad Request"),
		},
		{
			name: "Error closing response body",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("http: response body close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}
			result, err := r.handleCreate(context.Background(), tt.feed)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestFeedReconciler_handleDelete(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		expectedResult ctrl.Result
		expectedErr    bool
	}{
		{
			name: "Successful delete",
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
			name: "Feed name is empty",
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
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name: "Non-200 status code",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    true,
		},
		{
			name: "Error closing response body",
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

			result, err := r.handleDelete(context.Background(), tt.feed)

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestFeedReconciler_handleUpdate(t *testing.T) {
	tests := []struct {
		name           string
		feed           *newsaggregatorv1.Feed
		expectedResult ctrl.Result
		expectedErr    error
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
			expectedErr:    nil,
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
			expectedErr:    errors.New("json: error calling MarshalJSON"),
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
			expectedErr:    errors.New("http: invalid character"),
		},
		{
			name: "HTTP request execution error",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("dial tcp: lookup"),
		},
		{
			name: "Non-200 status code",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("Bad Request"),
		},
		{
			name: "Error closing response body",
			feed: &newsaggregatorv1.Feed{
				Spec: newsaggregatorv1.FeedSpec{
					Name: "Test Feed",
					Link: "http://example.com",
				},
			},
			expectedResult: ctrl.Result{},
			expectedErr:    errors.New("http: response body close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			result, err := r.handleUpdate(context.Background(), tt.feed)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func Test_initFeedStatus(t *testing.T) {
	tests := []struct {
		name      string
		feed      *newsaggregatorv1.Feed
		eventType string
		status    bool
		reason    string
		message   string
		expected  map[string]newsaggregatorv1.FeedConditions
	}{
		{
			name:      "Initializes the feed status correctly",
			feed:      &newsaggregatorv1.Feed{},
			eventType: "Ready",
			status:    true,
			reason:    "FeedCreated",
			message:   "Feed has been successfully created",
			expected: map[string]newsaggregatorv1.FeedConditions{
				"Ready": {
					Status:         true,
					Reason:         "FeedCreated",
					Message:        "Feed has been successfully created",
					LastUpdateTime: time.Now().Format(time.RFC3339),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			r.initFeedStatus(context.TODO(), tt.feed, tt.eventType, tt.status, tt.reason, tt.message)

			condition, exists := tt.feed.Status.Conditions[tt.eventType]
			assert.True(t, exists)
			assert.Equal(t, tt.expected[tt.eventType].Status, condition.Status)
			assert.Equal(t, tt.expected[tt.eventType].Reason, condition.Reason)
			assert.Equal(t, tt.expected[tt.eventType].Message, condition.Message)
		})
	}
}

func Test_updateFeedStatus(t *testing.T) {
	tests := []struct {
		name      string
		feed      *newsaggregatorv1.Feed
		eventType string
		status    bool
		reason    string
		message   string
		expected  map[string]newsaggregatorv1.FeedConditions
	}{
		{
			name: "Updates the feed status correctly",
			feed: &newsaggregatorv1.Feed{
				Status: newsaggregatorv1.FeedStatus{
					Conditions: map[string]newsaggregatorv1.FeedConditions{},
				},
			},
			eventType: "Ready",
			status:    true,
			reason:    "FeedUpdated",
			message:   "Feed status has been updated",
			expected: map[string]newsaggregatorv1.FeedConditions{
				"Ready": {
					Status:         true,
					Reason:         "FeedUpdated",
					Message:        "Feed status has been updated",
					LastUpdateTime: time.Now().Format(time.RFC3339),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FeedReconciler{}

			r.updateFeedStatus(context.TODO(), tt.feed, tt.eventType, tt.status, tt.reason, tt.message)

			condition, exists := tt.feed.Status.Conditions[tt.eventType]
			assert.True(t, exists)
			assert.Equal(t, tt.expected[tt.eventType].Status, condition.Status)
			assert.Equal(t, tt.expected[tt.eventType].Reason, condition.Reason)
			assert.Equal(t, tt.expected[tt.eventType].Message, condition.Message)
		})
	}
}
