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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
	}{
		{
			name: "Successful validation",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "UniqueFeedName",
					Link: "https://example.com/feed",
				},
			},
			expectedErr: false,
		},
		{
			name: "Validation failure due to invalid feed name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "SuperVeryLongNameThatIsNotValid",
					Link: "",
				},
			},
			expectedErr: true,
		},
		{
			name: "Validation failure due to invalid feed link",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "SuperVeryLongNameThatIsNotValid",
					Link: "",
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.feed.validateFeed()

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
