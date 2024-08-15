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

// HotNewsSpec defines the desired state of HotNews.
//
// This struct will be used to retrieve news by the criteria, specified here
// For example, we can specify keywords, date range, feeds and feed groups
// And then we will make requests to our news aggregator server with this parameters, and get the news
type HotNewsSpec struct {
	// Keywords is a comma-separated list of keywords which will be used to search news
	// +optional
	Keywords string `json:"keywords"`

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

	var articles []string

	for i := 0; i <= len(articlesTitles)-1 && i < r.Spec.SummaryConfig.TitlesCount-1; i++ {
		articles = append(articles, articlesTitles[i])
	}
	r.Status.ArticlesTitles = articles
}
