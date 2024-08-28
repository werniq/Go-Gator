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

package v1

import (
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

func TestFeed_validateFeed(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	existingFeedList := &FeedList{
		Items: []Feed{
			{
				Spec: FeedSpec{
					Name: "ExistingFeedName",
					Link: "https://example.com/feed",
				},
			},
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()

	var tests = []struct {
		name        string
		feed        *Feed
		expectedErr bool
		setup       func()
	}{
		{
			name: "Successful validation",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com",
				},
			},
			setup:       func() {},
			expectedErr: false,
		},
		{
			name: "Validation failure due to invalid feed link",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "SuperVeryLongNameThatIsNotValid",
					Link: "",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation failure due to duplicate link",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "Not-duplicate",
					Link: "https://example.com/feed",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation failure due to duplicate name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "ExistingFeedName",
					Link: "https://example.com/feed",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation failure due to invalid feed link v2",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "",
					Link: "ftp:/example.com",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "K8sClient List error",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "",
					Link: "",
				},
			},
			setup: func() {
				k8sClient = fake.NewClientBuilder().Build()
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := tt.feed.validateFeed()

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeed_ValidateCreate(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	_ = v12.AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&v12.ConfigMap{
			TypeMeta:   v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{},
			Data:       nil,
		}).
		Build()

	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
		setup   func()
	}{
		{
			name: "Validate delete",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com",
				},
				Status: FeedStatus{},
			},
			setup: func() {

			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate feed name",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com",
				},
				Status: FeedStatus{},
			},
			setup: func() {
				k8sClient = fake.NewClientBuilder().
					WithScheme(scheme).
					WithLists(&FeedList{
						Items: []Feed{
							{
								Spec: FeedSpec{
									Name: "UniqueFeedName",
									Link: "https://example.com",
								},
							},
						},
					}).Build()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := r.ValidateCreate()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.want, got, "ValidateCreate()")
		})
	}
}

func TestFeed_ValidateDelete(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}

	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	_ = v12.AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().
		WithScheme(scheme).
		WithLists(&HotNewsList{
			Items: []HotNews{
				{
					Spec: HotNewsSpec{
						Feeds:      []string{"sport"},
						FeedGroups: []string{"group1"},
					},
				},
			},
		}).WithObjects(&v12.ConfigMap{
		Data: map[string]string{"group1": "sport,news"},
	}).Build()

	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
		setup   func()
	}{
		{
			name: "Validate delete",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       FeedSpec{},
				Status:     FeedStatus{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Feed is used in hotnews",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "sport",
				},
				Status: FeedStatus{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Feed is used in feed groups",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "sport",
				},
				Status: FeedStatus{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := r.ValidateDelete()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.want, got, "ValidateDelete()")
		})
	}
}

func TestFeed_ValidateUpdate(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	_ = v12.AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&v12.ConfigMap{
			TypeMeta:   v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{},
			Data:       nil,
		}).
		Build()

	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
		setup   func()
	}{
		{
			name: "Validate delete",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com",
				},
				Status: FeedStatus{},
			},
			setup: func() {

			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate feed name",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com",
				},
				Status: FeedStatus{},
			},
			setup: func() {
				k8sClient = fake.NewClientBuilder().
					WithScheme(scheme).
					WithLists(&FeedList{
						Items: []Feed{
							{
								Spec: FeedSpec{
									Name: "UniqueFeedName",
									Link: "https://example.com",
								},
							},
						},
					}).Build()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := r.ValidateUpdate(r)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.want, got, "ValidateCreate()")
		})
	}
}

func TestIsFeedUsed(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	_ = v12.AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			&v12.ConfigMap{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "default",
					Name:      "feed-groups",
				},
				Data: map[string]string{"group1": "sport,news"},
			},
		).
		Build()

	feed := &Feed{}

	tests := []struct {
		name       string
		feedName   string
		hotNews    []HotNews
		configMaps []v12.ConfigMap
		setup      func()
		wantErr    bool
	}{
		{
			name:     "Feed used in HotNews",
			feedName: "sport",
			hotNews: []HotNews{
				{
					Spec: HotNewsSpec{
						Feeds: []string{"sport"},
					},
				},
			},
			configMaps: []v12.ConfigMap{
				{
					ObjectMeta: v1.ObjectMeta{
						Namespace: "default",
						Name:      "feed-groups",
					},
					Data: map[string]string{"group1": "sport,news"},
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name:     "Feed used in ConfigMap",
			feedName: "sport",
			hotNews: []HotNews{
				{
					Spec: HotNewsSpec{
						FeedGroups: []string{"group1"},
					},
				},
			},
			configMaps: []v12.ConfigMap{
				{
					ObjectMeta: v1.ObjectMeta{
						Namespace: "default",
						Name:      "feed-groups",
					},
					Data: map[string]string{"group1": "sport,news"},
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name:     "K8s Client Hot news List error",
			hotNews:  nil,
			feedName: "music",
			setup: func() {
				k8sClient = fake.NewClientBuilder().Build()
			},
			wantErr: true,
		},
		{
			name:     "K8s Client Config Map List error",
			hotNews:  nil,
			feedName: "music",
			setup: func() {
				scheme = runtime.NewScheme()
				_ = AddToScheme(scheme)
				k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if tt.configMaps != nil {
				if tt.hotNews != nil {
					k8sClient = fake.NewClientBuilder().
						WithScheme(scheme).
						WithLists(&v12.ConfigMapList{Items: tt.configMaps},
							&HotNewsList{Items: tt.hotNews},
						).
						Build()
				} else {
					k8sClient = fake.NewClientBuilder().
						WithScheme(scheme).
						WithLists(&v12.ConfigMapList{Items: tt.configMaps}).
						Build()
				}
			}

			err := feed.isFeedUsed(tt.feedName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedIsInFeedGroups(t *testing.T) {
	feed := &Feed{}

	tests := []struct {
		name       string
		feedGroups map[string]string
		feed       string
		expected   bool
	}{
		{
			name: "Feed is in a single group",
			feedGroups: map[string]string{
				"group1": "sport,news,music",
			},
			feed:     "sport",
			expected: true,
		},
		{
			name: "Feed is not in the groups",
			feedGroups: map[string]string{
				"group1": "news,music",
			},
			feed:     "sport",
			expected: false,
		},
		{
			name: "Feed is in multiple groups",
			feedGroups: map[string]string{
				"group1": "news,music",
				"group2": "sport,tech",
			},
			feed:     "sport",
			expected: true,
		},
		{
			name: "Feed is a substring but not a full match",
			feedGroups: map[string]string{
				"group1": "sports,news,music",
			},
			feed:     "sport",
			expected: false,
		},
		{
			name: "Empty feed groups",
			feedGroups: map[string]string{
				"group1": "",
			},
			feed:     "sport",
			expected: false,
		},
		{
			name: "Empty feed name",
			feedGroups: map[string]string{
				"group1": "sport,news,music",
			},
			feed:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := feed.feedIsInFeedGroups(tt.feedGroups, tt.feed)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFeed_checkLinkUniqueness(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}

	k8sClient = fake.NewClientBuilder().Build()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Error because hot news is not registered in schema",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       FeedSpec{},
				Status:     FeedStatus{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			_, err := r.checkLinkUniqueness()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeed_checkNameUniqueness(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}

	k8sClient = fake.NewClientBuilder().Build()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Error because hot news is not registered in schema",
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       FeedSpec{},
				Status:     FeedStatus{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			_, err := r.checkNameUniqueness()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFeeds_SetupWebhookWithManager(t *testing.T) {
	schema := runtime.NewScheme()
	assert.Nil(t, AddToScheme(schema))
	assert.Nil(t, v12.AddToScheme(schema))

	mgr, err := controllerruntime.NewManager(controllerruntime.GetConfigOrDie(), controllerruntime.Options{
		Scheme: schema,
	})
	assert.Nil(t, err)

	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	type args struct {
		mgr controllerruntime.Manager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful webhook setup",
			fields: fields{
				Spec:       FeedSpec{},
				Status:     FeedStatus{},
				ObjectMeta: v1.ObjectMeta{},
				TypeMeta:   v1.TypeMeta{},
			},
			args: args{
				mgr: mgr,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Feed{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}

			if tt.wantErr {
				assert.NotNil(t, r.SetupWebhookWithManager(tt.args.mgr))
			} else {
				assert.Nil(t, r.SetupWebhookWithManager(tt.args.mgr))
			}
		})
	}
}
