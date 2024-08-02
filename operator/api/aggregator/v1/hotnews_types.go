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
	// Keywords field contains a list of required keywords
	// +kubebuilder:validation:Required
	Keywords []string `json:"keywords"`

	// DateStart identifies minimal publication date from which news will be fetched
	DateStart string `json:"dateStart,omitempty"`

	// DateEnd identifies maximal publication date from which news will be fetched
	DateEnd string `json:"dateEnd,omitempty"`

	// Feeds are all feed names in the current namespace, if empty — will watch ALL available feed
	Feeds []string `json:"feeds"`

	// FeedGroups are all available sections of feeds from `feed-group-source` ConfigMap
	FeedGroups []string `json:"feedGroups"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// ArticlesCount is amount of articles by certain criteria
	ArticlesCount int `json:"articlesCount"`

	// NewsLink link to the news-aggregator HTTPs server to get all news by the criteria in JSON format
	NewsLink string `json:"newsLink"`

	// ArticlesTitles are first spec.summaryConfig.titlesCount article titles, sorted by feed name
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
