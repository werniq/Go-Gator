package v1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestFeed_SetCreatedCondition(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	type args struct {
		reason string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test SetCreatedCondition",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Namespace: "default",
					Name:      "test",
				},
				Spec:   FeedSpec{},
				Status: FeedStatus{},
			},
			args: args{
				reason: "FeedCreated",
			},
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
			r.SetCreatedCondition(tt.args.reason)
			assert.Equal(t, r.Status.Conditions[TypeFeedCreated].Reason, tt.args.reason)
		})
	}
}

func TestFeed_SetFailedCondition(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	type args struct {
		message string
		reason  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test SetCreatedCondition",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Namespace: "default",
					Name:      "test",
				},
				Spec:   FeedSpec{},
				Status: FeedStatus{},
			},
			args: args{
				reason:  "FeedFailedToCreate",
				message: "Failed to create feed",
			},
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
			r.SetFailedCondition(tt.args.message, tt.args.reason)
			assert.Equal(t, r.Status.Conditions[TypeFeedFailedToCreate].Reason, tt.args.reason)
		})
	}
}

func TestFeed_SetUpdatedCondition(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	type args struct {
		reason string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test SetCreatedCondition",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Namespace: "default",
					Name:      "test",
				},
				Spec:   FeedSpec{},
				Status: FeedStatus{},
			},
			args: args{
				reason: "FeedFailedToCreate",
			},
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
			r.SetUpdatedCondition(tt.args.reason)
			assert.Equal(t, r.Status.Conditions[TypeFeedUpdated].Reason, tt.args.reason)
		})
	}
}

func TestFeed_setCondition(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       FeedSpec
		Status     FeedStatus
	}
	type args struct {
		conditionType string
		status        bool
		reason        string
		message       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test SetCreatedCondition",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Namespace: "default",
					Name:      "test",
				},
				Spec:   FeedSpec{},
				Status: FeedStatus{},
			},
			args: args{
				conditionType: TypeFeedFailedToCreate,
				status:        false,
				message:       "Failed to create feed",
				reason:        "FeedFailedToCreate",
			},
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
			r.setCondition(tt.args.conditionType, tt.args.status, tt.args.reason, tt.args.message)
			assert.Equal(t, r.Status.Conditions[TypeFeedFailedToCreate].Reason, tt.args.reason)
		})
	}
}
