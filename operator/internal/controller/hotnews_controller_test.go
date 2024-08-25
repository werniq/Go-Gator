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
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"

	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

func TestHotNewsReconciler_Reconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"totalAmount": 2, "news": [{"title": "News 1"}, {"title": "News 2"}]}`))
	}))

	existingHotNewsList := newsaggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "feed-sample",
		},
		Spec: newsaggregatorv1.HotNewsSpec{
			Keywords:  []string{"keyword1,keyword2"},
			DateStart: "2024-08-12",
			DateEnd:   "2024-08-13",
			Feeds:     []string{"abc", "bbc"},
		},
	}

	existingConfigMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: newsaggregatorv1.FeedGroupsNamespace,
			Name:      newsaggregatorv1.FeedGroupsConfigMapName,
		},
		Data: map[string]string{
			"sport":   "washingtontimes",
			"politic": "abc,bbc",
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).
		WithObjects(&existingHotNewsList, &existingConfigMap).
		Build()

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
				Client: k8sClient,
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
				Client: k8sClient,
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
				Client: k8sClient,
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
				Client:    k8sClient,
				Scheme:    tt.fields.Scheme,
				serverUrl: mockServer.URL,
			}

			got, err := r.Reconcile(context.TODO(), tt.args.req)
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
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)
	serverNewsEndpoint := "https://go-gator-svc.go-gator.svc.cluster.local:443/news"
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
					Keywords:  []string{"bitcoin"},
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
					Keywords:  []string{"bitcoin"},
					Feeds:     []string{"abc", "bbc"},
					DateStart: "2024-08-05",
				},
			},
			want:    serverNewsEndpoint + "?keywords=bitcoin&sources=abc,bbc&dateFrom=2024-08-05",
			wantErr: false,
		},
		{
			name: "Valid request with keywords, feedGroups, and start date only",
			fields: fields{
				Client: fake.NewClientBuilder().WithScheme(scheme).
					WithObjects(&v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: newsaggregatorv1.FeedGroupsNamespace,
							Name:      newsaggregatorv1.FeedGroupsConfigMapName,
						},
						Data: map[string]string{"sport": "abc,bbc"},
					}).
					Build(),
			},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords:   []string{"bitcoin"},
					FeedGroups: []string{"sport"},
					DateStart:  "2024-08-05",
				},
			},
			want:    serverNewsEndpoint + "?keywords=bitcoin&sources=abc,bbc&dateFrom=2024-08-05",
			wantErr: false,
		},
		{
			name:   "Valid request with keywords and feeds only",
			fields: fields{},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					Keywords: []string{"bitcoin"},
					Feeds:    []string{"abc", "bbc"},
				},
			},
			want:    serverNewsEndpoint + "?keywords=bitcoin&sources=abc,bbc",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNewsReconciler{
				Client:    tt.fields.Client,
				Scheme:    tt.fields.Scheme,
				serverUrl: serverNewsEndpoint,
			}
			got, err := r.constructRequestUrl(context.Background(), tt.args.spec)
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
	_ = v1.AddToScheme(scheme)

	existingConfigMap := v1.ConfigMapList{
		Items: []v1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: newsaggregatorv1.FeedGroupsNamespace,
					Name:      newsaggregatorv1.FeedGroupsConfigMapName,
				},
				Data: map[string]string{
					"sport":   "washingtontimes",
					"politic": "abc,bbc",
				},
			},
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).
		WithLists(&existingConfigMap).
		Build()

	errorClient := &errorReturningClient{
		Client: fakeClient,
	}

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
		want    v1.ConfigMapList
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
			},
		},
		{
			name: "Error during listing of ConfigMaps",
			fields: fields{
				Client: errorClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    v1.ConfigMapList{},
			wantErr: true,
			setup: func() {
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

			got, err := r.getFeedGroups(context.Background())
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

// Custom client that returns an error during the List operation
type errorReturningClient struct {
	client.Client
}

func (c *errorReturningClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return fmt.Errorf("simulated error during listing of ConfigMaps")
}

func TestHotNewsReconciler_processHotNews(t *testing.T) {
	serverNewsEndpoint := "https://go-gator-svc.go-gator.svc.cluster.local:443/news"
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	existingFeedList := &newsaggregatorv1.HotNewsList{
		Items: []newsaggregatorv1.HotNews{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "feed-sample",
				},
				Spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  []string{"keyword1,keyword2"},
					DateStart: "2024-08-12",
					DateEnd:   "2024-08-13",
				},
			},
		},
	}

	existingConfigMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: newsaggregatorv1.FeedGroupsNamespace,
			Name:      newsaggregatorv1.FeedGroupsConfigMapName,
		},
		Data: map[string]string{
			"sport":   "washingtontimes",
			"politic": "abc,bbc",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).
		WithLists(existingFeedList).
		WithObjects(&existingConfigMap).
		Build()

	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx     context.Context
		hotNews *newsaggregatorv1.HotNews
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockServer *httptest.Server
		wantErr    bool
	}{
		{
			name: "Successful execution",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
				hotNews: &newsaggregatorv1.HotNews{
					Spec: newsaggregatorv1.HotNewsSpec{
						Keywords:  []string{"bitcoin"},
						DateStart: "2024-08-05",
						DateEnd:   "2024-08-06",
						Feeds:     []string{"abc", "bbc"},
					},
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"totalAmount": 2, "news": [{"title": "News 1"}, {"title": "News 2"}]}`))
			})),
			wantErr: false,
		},
		{
			name: "Failed to construct request URL",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
				hotNews: &newsaggregatorv1.HotNews{
					Spec: newsaggregatorv1.HotNewsSpec{
						Keywords: []string{"bitcoin"},
						Feeds:    []string{"abc", "bbc"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed to create HTTP request",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
				hotNews: &newsaggregatorv1.HotNews{
					Spec: newsaggregatorv1.HotNewsSpec{
						Keywords:  []string{"bitcoin"},
						DateStart: "2024-08-05",
						DateEnd:   "2024-08-06",
						Feeds:     []string{"abc", "bbc"},
					},
				},
			},
			mockServer: nil,
			wantErr:    true,
		},
		{
			name: "HTTP request failed",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
				hotNews: &newsaggregatorv1.HotNews{
					Spec: newsaggregatorv1.HotNewsSpec{
						Keywords:  []string{"bitcoin"},
						DateStart: "2024-08-05",
						DateEnd:   "2024-08-06",
						Feeds:     []string{"abc", "bbc"},
					},
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			})),
			wantErr: true,
		},
		{
			name: "Invalid response body JSON",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				ctx: context.TODO(),
				hotNews: &newsaggregatorv1.HotNews{
					Spec: newsaggregatorv1.HotNewsSpec{
						Keywords:  []string{"bitcoin"},
						DateStart: "2024-08-05",
						DateEnd:   "2024-08-06",
						Feeds:     []string{"abc", "bbc"},
					},
				},
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{invalid json}`))
			})),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockServer != nil {
				defer tt.mockServer.Close()
				serverNewsEndpoint = tt.mockServer.URL
			}

			r := &HotNewsReconciler{
				Client:    tt.fields.Client,
				Scheme:    tt.fields.Scheme,
				serverUrl: serverNewsEndpoint,
			}
			if err := r.processHotNews(tt.args.ctx, tt.args.hotNews); (err != nil) != tt.wantErr {
				t.Errorf("processHotNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNewsReconciler_processFeedGroups(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	existingFeedList := &newsaggregatorv1.HotNewsList{
		Items: []newsaggregatorv1.HotNews{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "feed-sample",
				},
				Spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  []string{"keyword1,keyword2"},
					DateStart: "2024-08-12",
					DateEnd:   "2024-08-13",
				},
			},
		},
	}

	existingConfigMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: newsaggregatorv1.FeedGroupsNamespace,
			Name:      newsaggregatorv1.FeedGroupsConfigMapName,
		},
		Data: map[string]string{
			"sport":   "washingtontimes",
			"politic": "abc,bbc",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).
		WithLists(existingFeedList).
		WithObjects(&existingConfigMap).
		Build()

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
		setup   func() *v1.ConfigMap
		want    string
		wantErr bool
	}{
		{
			name: "Successful processing with valid feed groups",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					FeedGroups: []string{"sport", "politic"},
				},
			},
			setup: func() *v1.ConfigMap {
				return &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      newsaggregatorv1.FeedGroupsConfigMapName,
						Namespace: newsaggregatorv1.FeedGroupsNamespace,
					},
					Data: map[string]string{
						"sport":   "washingtontimes",
						"politic": "abc,bbc",
					},
				}
			},
			want:    "washingtontimes,abc,bbc",
			wantErr: false,
		},
		{
			name: "Feed group not found in ConfigMap",
			fields: fields{
				Client: fakeClient,
				Scheme: scheme,
			},
			args: args{
				spec: newsaggregatorv1.HotNewsSpec{
					FeedGroups: []string{"nonexistent"},
				},
			},
			setup: func() *v1.ConfigMap {
				return &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      newsaggregatorv1.FeedGroupsConfigMapName,
						Namespace: newsaggregatorv1.FeedGroupsNamespace,
					},
					Data: map[string]string{
						"sport": "washingtontimes",
					},
				}
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				configMap := tt.setup()
				if configMap != nil {
					tt.fields.Client = fake.NewClientBuilder().
						WithScheme(scheme).
						WithObjects(configMap).
						Build()
				}
			}

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
				t.Errorf("processFeedGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHotNewsReconciler_SetupWithManager(t *testing.T) {
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
			r := &HotNewsReconciler{
				serverUrl: tt.fields.serverUrl,
				Client:    tt.fields.Client,
				Scheme:    tt.fields.Scheme,
			}
			err := r.SetupWithManager(tt.args.mgr, tt.args.serverUrl)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.args.serverUrl, r.serverUrl)
		})
	}
}
