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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// TypeHotNewsCreated represents the Created condition type
	TypeHotNewsCreated = "Created"

	// TypeHotNewsUpdated represents the reason for created condition
	TypeHotNewsUpdated = "Updated"

	// HotNewsSuccessfullyCreated represents the reason for created condition
	HotNewsSuccessfullyCreated = "HotNews was successfully created"

	// HotNewsError indicates that there were an error during Reconciliation of hot news object
	HotNewsError = "Error during processing of hot news creation"

	// StatusError is set to false, which says that there was an error and Reconciliation was not
	// completed successfully
	StatusError = false

	// StatusSuccess indicates that HotNews object was successfully created or updated
	StatusSuccess = true

	// hotNewsStatusConditionsCapacity is a capacity of hot news status conditions map
	// It is defaulted to 3, since we have 3 conditions: Created, Updated, FailedToCreate.
	// Condition Deleted is not included, since it is not used in the current implementation
	hotNewsStatusConditionsCapacity = 3
)

// HotNewsSpec defines the desired state of HotNews.
//
// This struct will be used to retrieve news by the criteria, specified here
// For example, we can specify keywords, date range, feeds and feed groups
// And then we will make requests to our news aggregator server with this parameters, and get the news
type HotNewsSpec struct {
	// Keywords is a comma-separated list of keywords which will be used to search news
	// +kubebuilder:validation:Required
	Keywords []string `json:"keywords"`

	// DateStart is a news starting date in format "YYYY-MM-DD", can be empty
	// +optional
	DateStart string `json:"dateStart,omitempty"`

	// DateEnd is a news final date in format "YYYY-MM-DD", can be empty
	DateEnd string `json:"dateEnd,omitempty"`

	// Feeds is a list of Feeds CRD, which will be used to subscribe to news
	// +optional
	Feeds []string `json:"feeds,omitempty"`

	// FeedGroups are available sections of feeds from `hotNew-group-source` ConfigMap
	// +optional
	FeedGroups []string `json:"feedGroups,omitempty"`

	// SummaryConfig summary of observed hot news
	// +optional
	SummaryConfig SummaryConfig `json:"summaryConfig,omitempty"`
}

// SummaryConfig struct defines the configuration for the summary of hot news
// It stores the number of titles to show and store in HotNewsStatus.ArticlesTitles
type SummaryConfig struct {
	// TitlesCount is a number of titles to show in the summary
	TitlesCount int `json:"titlesCount"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// ArticlesCount displays total amount of news by the criteria
	ArticlesCount int `json:"articlesCount"`

	// NewsLink is a link which will be constructed to get all news by the certain criteria
	NewsLink string `json:"newsLink"`

	// ArticlesTitles contains a list of titles of first 10 articles
	ArticlesTitles []string `json:"articlesTitles"`

	Conditions map[string]HotNewsConditions `json:"conditions,omitempty"`
}

type HotNewsConditions struct {
	// Status field is a boolean that represents the status of the condition
	// A value of true typically indicates the condition is met, while
	// false indicates it is not.
	Status bool `json:"status"`

	// Reason field is a string which is populated if status is false
	// It explains the reason for the current status.
	Reason string `json:"reason"`

	// Message field is a string which is populated if status is false
	// It provides additional details or a message about the condition.
	Message string `json:"message"`

	// LastUpdateTime is a time when an object changes its state
	// This timestamp indicates the last time the condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}

// GetFeedGroupNames returns all config maps which contain hotNew groups names
func (r *HotNews) GetFeedGroupNames(ctx context.Context) ([]string, error) {
	s, err := labels.NewRequirement(FeedGroupLabel, selection.Exists, nil)
	if err != nil {
		return nil, err
	}

	var configMaps v1.ConfigMapList
	err = k8sClient.List(ctx, &configMaps, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*s),
		Namespace:     r.Namespace,
	})
	if err != nil {
		return nil, err
	}

	var feedGroups []string
	for _, configMap := range configMaps.Items {
		for _, source := range r.Spec.FeedGroups {
			if _, exists := configMap.Data[source]; exists {
				feedGroups = append(feedGroups, source)
			}
		}
	}

	return feedGroups, nil
}

// SetCondition initializes status.Conditions if they are empty.
// In this case, capacity of Conditions mapping is 3, since we support 3 conditions:
// 1. Created
// 2. Updated
// 3. Error
// If Conditions mapping already is not nil, this functions creates or updates condition with given
// condition type, reason, and message.
// LastUpdateTime is set by default to metav1.Now()
func (r *HotNews) SetCondition(conditionType string, status bool, reason, message string) {
	if r.Status.Conditions == nil {
		r.Status.Conditions = make(map[string]HotNewsConditions, hotNewsStatusConditionsCapacity)
	}

	r.Status.Conditions[conditionType] = HotNewsConditions{
		Status:         status,
		Reason:         reason,
		Message:        message,
		LastUpdateTime: metav1.Now(),
	}
}

// SetStatus func initializes HotNews.Status object with the provided data
func (r *HotNews) SetStatus(articlesCount int, requestUrl string, articlesTitles []string) {
	r.Status.ArticlesCount = articlesCount
	r.Status.NewsLink = requestUrl

	var articles []string

	for i := 0; i <= len(articlesTitles)-1 && i < r.Spec.SummaryConfig.TitlesCount; i++ {
		articles = append(articles, articlesTitles[i])
	}
	r.Status.ArticlesTitles = articles
}
