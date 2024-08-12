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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec defines the desired state of HotNews
type HotNewsSpec struct {
	// Keywords is a comma-separated list of keywords which will be used to search news
	// +kubebuilder:validation:Required
	Keywords string `json:"keywords"`

	// DateStart is a news starting date in format "YYYY-MM-DD", can be empty
	DateStart string `json:"dateStart"`

	// DateEnd is a news final date in format "YYYY-MM-DD", can be empty
	DateEnd string `json:"dateEnd"`

	// Feeds is a list of Feeds CRD, which will be used to subscribe to news
	Feeds []string `json:"feeds"`

	// FeedGroups are available sections of feeds from `feed-group-source` ConfigMap
	FeedGroups []string `json:"feedGroups"`

	// SummaryConfig summary of observed hot news
	SummaryConfig SummaryConfig `json:"summaryConfig"`
}

// SummaryConfig struct
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
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}

// InitHotNewsStatus func initializes HotNews.Status object with the provided data
func (r *HotNews) InitHotNewsStatus(articlesCount int, requestUrl string, articlesTitles []string) {
	r.Status.ArticlesCount = articlesCount
	r.Status.NewsLink = requestUrl
	r.Status.ArticlesTitles = articlesTitles
}
