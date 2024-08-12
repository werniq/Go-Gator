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
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

func TestValidateFeed(t *testing.T) {
	var tests = []struct {
		name           string
		feed           *Feed
		mockResponse   string
		mockError      error
		expectedError  string
		shouldPanic    bool
		validationErr  error
		k8sClientSetup func(k8sClient *fake.Clientset)
	}{
		{
			name: "Successful validation",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com/feed",
				},
			},
			mockResponse:  `{"items":[]}`,
			expectedError: "",
			validationErr: nil,
		},
		{
			name: "Validation failure due to invalid feed name or link",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "",
					Link: "",
				},
			},
			mockResponse:  "",
			expectedError: "validation error",
			validationErr: errors.New("validation error"),
		},
		{
			name: "Feed name already exists",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "ExistingFeedName",
					Link: "https://example.com/feed",
				},
			},
			mockResponse:  `{"items":[{"spec":{"name":"ExistingFeedName"}}]}`,
			expectedError: "name must be unique",
		},
		{
			name: "K8s client error",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com/feed",
				},
			},
			mockResponse:  "",
			expectedError: "error creating Kubernetes client",
		},
		{
			name: "Unmarshal error",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com/feed",
				},
			},
			mockResponse:  "invalid json",
			expectedError: "invalid character",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, err := tt.feed.validateFeed()

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.Nil(t, err)
				assert.Nil(t, warnings)
			}
		})
	}
}

func TestFeed_validateFeed(t *testing.T) {
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful validation",
			fields: fields{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com/feed",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			fields: fields{
				Spec: FeedSpec{
					Name: "SuperVeryLongNameThatIsNotValid",
					Link: "https://example.com/feed",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			fields: fields{
				Spec: FeedSpec{
					Name: "Normal Nam",
					Link: "wrong.com/feed",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
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
			got, err := r.validateFeed()
			if !tt.wantErr(t, err, fmt.Sprintf("validateFeed()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "validateFeed()")
		})
	}
}
