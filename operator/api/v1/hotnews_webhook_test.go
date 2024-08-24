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
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestFeed_validateHotNews(t *testing.T) {
	var err error
	clientset, err = kubernetes.NewForConfig(c)
	assert.Nil(t, err)

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: FeedGroupsNamespace,
		},
	}, v12.CreateOptions{
		FieldManager: "test",
	})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			err = nil
		}
	}
	assert.Nil(t, err)

	_, err = clientset.CoreV1().ConfigMaps(FeedGroupsNamespace).Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Name: FeedGroupsConfigMapName,
		},
	}, v12.CreateOptions{
		FieldManager: "test",
	})
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
					Keywords:  []string{"test"},
					DateStart: "2021-01-01",
					DateEnd:   "2021-01-02",
					Feeds:     []string{"feed1"},
				},
			},
			expectedErr: false,
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
			expectedErr: true,
		},
		{
			name: "Validation failure because of empty feeds and feedGroups",
			hotNew: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"test"},
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
					Keywords:   []string{"test"},
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
