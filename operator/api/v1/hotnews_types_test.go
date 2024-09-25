package v1

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"testing"
)

func TestGetFeedGroupNames(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v12.AddToScheme(scheme)

	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	var res []string
	testCases := []struct {
		name           string
		k8sObjects     []runtime.Object
		expectedGroups []string
		expectError    bool
		setup          func()
		client         client.Client
		feedGroups     []string
	}{
		{
			name:           "No config maps exist",
			k8sObjects:     []runtime.Object{},
			feedGroups:     []string{"group1", "group2"},
			expectedGroups: res,
			client:         c,
			expectError:    false,
		},
		{
			name: "Config map with no matching feed groups",
			k8sObjects: []runtime.Object{
				&v12.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:   "config2",
						Labels: map[string]string{FeedGroupLabel: "true"},
					},
					Data: map[string]string{
						"othergroup": "some-data",
					},
				},
			},
			feedGroups:     []string{"group1", "group2"},
			expectedGroups: res,
			client:         c,
			setup:          func() {},
			expectError:    false,
		},
		{
			name:       "Config map listing error",
			k8sObjects: nil,
			feedGroups: []string{},
			setup: func() {
				scheme = runtime.NewScheme()
			},
			expectedGroups: res,
			client: fake.NewClientBuilder().WithInterceptorFuncs(
				interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
					return errors.New("error")
				}},
			).Build(),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient = tc.client

			hotNews := &HotNews{
				Spec: HotNewsSpec{
					FeedGroups: tc.feedGroups,
				},
			}

			result, err := hotNews.GetFeedGroupNames(context.TODO())

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedGroups, result)
		})
	}
}

func TestHotNews_InitHotNewsStatus(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
	}
	type args struct {
		articlesCount  int
		requestUrl     string
		articlesTitles []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test valid execution",
			args: args{
				articlesCount:  10,
				requestUrl:     "http://test.com",
				articlesTitles: []string{"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9", "test10"},
			},
			fields: fields{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec: HotNewsSpec{
					SummaryConfig: SummaryConfig{TitlesCount: 5},
				},
				Status: HotNewsStatus{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNews{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			r.SetStatus(tt.args.articlesCount, tt.args.requestUrl, tt.args.articlesTitles)
			assert.Equal(t, r.Status.ArticlesCount, tt.args.articlesCount)
		})
	}
}
