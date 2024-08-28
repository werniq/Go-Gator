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
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

func TestFeed_validateHotNews(t *testing.T) {
	k8sClient = fake.NewClientBuilder().
		WithObjects(&v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Namespace: FeedGroupsNamespace,
				Name:      FeedGroupsConfigMapName,
			},
			Data: map[string]string{"sport": "abc"},
		}).
		Build()

	var tests = []struct {
		name        string
		hotNew      *HotNews
		expectedErr bool
		setup       func()
	}{
		{
			name: "Successful validation",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
					Feeds:     []string{"feed1"},
				},
			},
			setup:       func() {},
			expectedErr: false,
		},
		{
			name: "Validation failure due to empty feeds and feed groups",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation failure due to invalid hotNew date range",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-03",
					DateEnd:   "2021-01-02",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation failure due to invalid dates",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "ABCC-AA-BB",
					DateEnd:   "BBCA-AA-BB",
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "Validation with feed groups",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:   []string{"test"},
					DateStart:  "ABCC-AA-BB",
					FeedGroups: []string{"sport"},
				},
			},
			setup:       func() {},
			expectedErr: true,
		},
		{
			name: "K8s client not est",
			hotNew: &HotNews{
				ObjectMeta: v12.ObjectMeta{
					Namespace: "non-eadssxistent",
				},
				Spec: HotNewsSpec{
					Keywords:   []string{"test"},
					DateStart:  "ABCC-AA-BB",
					FeedGroups: []string{"sport"},
				},
			},
			setup: func() {
				k8sClient = fake.NewClientBuilder().WithScheme(nil).Build()
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := tt.hotNew.validateHotNews()

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestHotNews_ValidateCreate(t *testing.T) {
	k8sClient = fake.NewClientBuilder().WithObjects(&v1.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Namespace: FeedGroupsNamespace,
			Name:      FeedGroupsConfigMapName,
		},
		Data: map[string]string{
			"sport": "abc",
		},
	}).Build()
	type fields struct {
		TypeMeta   v12.TypeMeta
		ObjectMeta v12.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Successful validation",
			fields: fields{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
					Feeds:     []string{"abc"},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Successful validation with feed groups",
			fields: fields{
				Spec: HotNewsSpec{
					Keywords:   []string{"test"},
					DateStart:  "2021-01-01",
					DateEnd:    "2021-01-02",
					FeedGroups: []string{"sport"},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Not Successful validation with feed groups",
			fields: fields{
				Spec: HotNewsSpec{
					Keywords:   []string{"test"},
					DateStart:  "2021-01-01",
					DateEnd:    "2021-01-02",
					FeedGroups: []string{"abbbc"},
				},
			},
			want:    nil,
			wantErr: true,
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

func TestHotNews_ValidateDelete(t *testing.T) {
	type fields struct {
		TypeMeta   v12.TypeMeta
		ObjectMeta v12.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name:    "Successful validation",
			wantErr: false,
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

func TestHotNews_ValidateUpdate(t *testing.T) {
	k8sClient = fake.NewClientBuilder().WithObjects(&v1.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Namespace: FeedGroupsNamespace,
			Name:      FeedGroupsConfigMapName,
		},
		Data: nil}).Build()

	type fields struct {
		TypeMeta   v12.TypeMeta
		ObjectMeta v12.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
	}
	type args struct {
		old runtime.Object
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Successful validation",
			fields: fields{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
					Feeds:     []string{"abc"},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Successful validation",
			fields: fields{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
				},
			},
			want:    nil,
			wantErr: true,
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
			got, err := r.ValidateUpdate(tt.args.old)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.want, got, "ValidateUpdate(%v)", tt.args.old)
		})
	}
}

func TestHotNews_Default(t *testing.T) {
	type fields struct {
		TypeMeta   v12.TypeMeta
		ObjectMeta v12.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Successful defaulting",
			fields: fields{
				Spec:       HotNewsSpec{},
				Status:     HotNewsStatus{},
				ObjectMeta: v12.ObjectMeta{},
				TypeMeta:   v12.TypeMeta{},
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
			r.Default()
		})
	}
}

func TestHotNews_SetupWebhookWithManager(t *testing.T) {
	schema := runtime.NewScheme()
	assert.Nil(t, AddToScheme(schema))
	assert.Nil(t, v1.AddToScheme(schema))

	mgr, err := controllerruntime.NewManager(controllerruntime.GetConfigOrDie(), controllerruntime.Options{
		Scheme: schema,
	})
	assert.Nil(t, err)

	type fields struct {
		TypeMeta   v12.TypeMeta
		ObjectMeta v12.ObjectMeta
		Spec       HotNewsSpec
		Status     HotNewsStatus
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
				Spec:       HotNewsSpec{},
				Status:     HotNewsStatus{},
				ObjectMeta: v12.ObjectMeta{},
				TypeMeta:   v12.TypeMeta{},
			},
			args: args{
				mgr: mgr,
			},
			wantErr: false,
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

			if tt.wantErr {
				assert.NotNil(t, r.SetupWebhookWithManager(tt.args.mgr))
			} else {
				assert.Nil(t, r.SetupWebhookWithManager(tt.args.mgr))
			}
		})
	}
}
