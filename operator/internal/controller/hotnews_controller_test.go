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
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"

	newsaggregatorv1 "teamdev.com/go-gator-operator/api/v1"
)

func TestHotNewsReconciler_Reconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)

	existingFeedList := &newsaggregatorv1.HotNewsList{
		Items: []newsaggregatorv1.HotNews{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "feed-sample",
				},
				Spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  "keyword1,keyword2",
					DateStart: "2024-08-12",
					DateEnd:   "2024-08-13",
				},
			},
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).
		WithLists(existingFeedList).Build()

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
			name: "Successful Reconcile call",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "feed-sample",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "Failed because feed not found (invalid name)",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "feed-not-found",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: true,
		},
		{
			name: "Failed because feed not found (invalid namespace)",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "non-existent-namespace",
						Name:      "feed-not-found",
					},
				},
			},
			want:    controllerruntime.Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
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

func TestHotNewsReconciler_constructRequestUrl(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		spec newsaggregatorv1.HotNewsSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "Valid request with keywords, feeds, and date range",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  "bitcoin",
					Feeds:     []string{"abc", "bbc"},
					DateStart: "2024-08-05",
					DateEnd:   "2024-08-06",
				},
			},
			want:    serverNewsEndpoint + "?keywords=bitcoin&sources=abc,bbc&dateFrom=2024-08-05&dateEnd=2024-08-06",
			wantErr: false,
		},
		{
			name:   "Valid request with keywords, feeds, and start date only",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  "bitcoin",
					Feeds:     []string{"abc", "bbc"},
					DateStart: "2024-08-05",
				},
			},
			want:    serverNewsEndpoint + "news?keywords=bitcoin&sources=abc,bbc&dateFrom=2024-08-05",
			wantErr: false,
		},
		{
			name:   "Valid request with keywords and feeds only",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords: "bitcoin",
					Feeds:    []string{"abc", "bbc"},
				},
			},
			want:    serverNewsEndpoint + "news?keywords=bitcoin&sources=abc,bbc",
			wantErr: false,
		},
		{
			name:   "Valid request with keywords and feed groups",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords:   "bitcoin",
					FeedGroups: []string{"group1", "group2"},
				},
			},
			want:    serverNewsEndpoint + "news?keywords=bitcoin&sources=group1,group2",
			wantErr: false,
		},
		{
			name:   "Feeds and feed groups both present, should prefer feed groups",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords:   "bitcoin",
					Feeds:      []string{"abc", "bbc"},
					FeedGroups: []string{"group1", "group2"},
				},
			},
			want:    serverNewsEndpoint + "news?keywords=bitcoin&sources=group1,group2",
			wantErr: false,
		},
		{
			name:   "Feeds present but empty, should return error",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords: "bitcoin",
					Feeds:    []string{},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:   "Missing feeds should return an error",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords: "bitcoin",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.constructRequestUrl(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("constructRequestUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("constructRequestUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHotNewsReconciler_getFeedGroups(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)

	existingFeedList := &newsaggregatorv1.HotNewsList{
		Items: []newsaggregatorv1.HotNews{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "feed-sample",
				},
				Spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  "keyword1,keyword2",
					DateStart: "2024-08-12",
					DateEnd:   "2024-08-13",
				},
			},
		},
	}

	// Define a ConfigMap that would be returned in the successful test cases
	existingConfigMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feed-group-source",
			Namespace: "operator-system",
		},
		Data: map[string]string{
			"sport":   "washingtontimes",
			"politic": "abc,bbc",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()

	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.ConfigMap
		wantErr bool
		setup   func()
	}{
		{
			name: "Successful retrieval of ConfigMap",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    existingConfigMap,
			wantErr: false,
			setup: func() {
				// Add ConfigMap to fake client
				err := fakeClient.Create(context.TODO(), existingConfigMap)
				assert.Nil(t, err)
			},
		},
		{
			name: "ConfigMap not found",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				// Ensure the ConfigMap does not exist in the fake client
				err := fakeClient.Delete(context.TODO(), existingConfigMap)
				assert.Nil(t, err)
			},
		},
		{
			name: "ConfigMap retrieved with no data",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "feed-group-source",
					Namespace: "operator-system",
				},
				Data: map[string]string{},
			},
			wantErr: false,
			setup: func() {
				// Create an empty ConfigMap in the fake client
				emptyConfigMap := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "feed-group-source",
						Namespace: "operator-system",
					},
					Data: map[string]string{},
				}
				err := fakeClient.Create(context.TODO(), emptyConfigMap)
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.getFeedGroups(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFeedGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFeedGroups() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHotNewsReconciler_handleCreate(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx     context.Context
		hotNews *newsaggregatorv1.HotNews
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.handleCreate(tt.args.ctx, tt.args.hotNews); (err != nil) != tt.wantErr {
				t.Errorf("handleCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNewsReconciler_handleDelete(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx     context.Context
		hotNews *newsaggregatorv1.HotNews
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.handleDelete(tt.args.ctx, tt.args.hotNews); (err != nil) != tt.wantErr {
				t.Errorf("handleDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNewsReconciler_handleUpdate(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx     context.Context
		hotNews *newsaggregatorv1.HotNews
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.handleUpdate(tt.args.ctx, tt.args.hotNews); (err != nil) != tt.wantErr {
				t.Errorf("handleUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNewsReconciler_processFeedGroups(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		spec newsaggregatorv1.HotNewsSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.processFeedGroups(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("processFeedGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processFeedGroups() got = %v, want %v", got, tt.want)
			}
		})
	}
}
