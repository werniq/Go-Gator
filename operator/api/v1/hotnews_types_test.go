package v1

import (
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestGetFeedGroupNames(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v12.AddToScheme(scheme)

	var res []string
	testCases := []struct {
		name           string
		input          v12.ConfigMapList
		k8sObjects     []runtime.Object
		expectedGroups []string
		expectError    bool
		setup          func()
		feedGroups     []string
	}{
		{
			name:           "No config maps exist",
			input:          v12.ConfigMapList{},
			k8sObjects:     []runtime.Object{},
			feedGroups:     []string{"group1", "group2"},
			expectedGroups: res,
			expectError:    false,
		},
		{
			name: "Config map with no matching feed groups",
			input: v12.ConfigMapList{
				Items: []v12.ConfigMap{
					{
						ObjectMeta: v1.ObjectMeta{
							Name:   "config2",
							Labels: map[string]string{FeedGroupLabel: "true"},
						},
						Data: map[string]string{
							"othergroup": "some-data",
						},
					},
				},
			},
			feedGroups:     []string{"group1", "group2"},
			expectedGroups: res,
			setup:          func() {},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			hotNews := &HotNews{
				Spec: HotNewsSpec{
					FeedGroups: tc.feedGroups,
				},
			}

			result := hotNews.GetFeedGroupNames(tc.input)
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
