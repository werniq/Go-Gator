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
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestFeed_validateHotNews(t *testing.T) {
	var err error
	clientset, err = kubernetes.NewForConfig(c)
	assert.Nil(t, err)
	var tests = []struct {
		name        string
		hotNew      *HotNews
		expectedErr bool
	}{
		{
			name: "Successful validation",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:   "test",
					DateStart:  "2021-01-01",
					DateEnd:    "2021-01-02",
					FeedGroups: []string{"sport"},
				},
			},
			expectedErr: false,
		},
		{
			name: "Validation failure due to invalid hotNew date range",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  "test",
					DateStart: "2021-01-03",
					DateEnd:   "2021-01-02",
				},
			},
			expectedErr: true,
		},
		{
			name: "Validation failure because of empty feeds and feedGroups",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  "test",
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
				},
			},
			expectedErr: true,
		},
		{
			name: "Validation failure because of empty feeds and feedGroups",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:   "test",
					DateStart:  "2021-01-01",
					DateEnd:    "2021-01-02",
					FeedGroups: []string{"non-existing-feed"},
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hotNew.validateHotNews()

			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
