package v1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

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
			r.InitHotNewsStatus(tt.args.articlesCount, tt.args.requestUrl, tt.args.articlesTitles)
			assert.Equal(t, r.Status.ArticlesCount, tt.args.articlesCount)
		})
	}
}
