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
	// Keywords is a list of specific keywords
	Keywords string `json:"keywords"`

	// DateStart news starting date, can be empty
	DateStart string `json:"dateStart"`

	// DateEnd news ending date, can be empty
	DateEnd string `json:"dateEnd"`

	// Feeds is a list of Feed CRD, which will be used to subscribe to news
	Feeds []string `json:"feeds"`

	// FeedGroups are available sections of feeds from `feed-group-source` ConfigMap
	FeedGroups []string `json:"feedGroups"`

	// SummaryConfig summary of observed hot news
	SummaryConfig SummaryConfig `json:"summaryConfig"`
}

// SummaryConfig struct
type SummaryConfig struct {
	// ArticlesCount displays total amount of news by the criteria
	ArticlesCount int `json:"articlesCount"`

	// NewsLink is a link which will be constructed to get all news by the certain criteria
	NewsLink string `json:"newsLink"`

	// ArticlesTitles contains a list of titles of first 10 articles
	ArticlesTitles []string `json:"articlesTitles"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
