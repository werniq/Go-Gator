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

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	// Name field is a string that represents the name of the feed
	Name string `json:"name,omitempty"`

	// Link field is a string that represents the URL of the feed
	Link string `json:"link,omitempty"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	// Conditions field is a list of conditions that the feed can have
	Conditions []FeedConditions `json:"conditions,omitempty"`
}

// FeedConditions are the cond
type FeedConditions struct {
	// Type field is a string that represents the type of the condition
	Type string `json:"type"`

	// Status field is a boolean that represents the status of the condition
	Status bool `json:"status"`

	// Reason field is a string which is populated if status is false
	Reason string `json:"reason"`

	// Message field is a string which is populated if status is false
	Message string `json:"message"`

	// LastUpdateTime is a time when an object changes it's state
	LastUpdateTime string `json:"lastUpdateTime"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Feed is the Schema for the feeds API
type Feed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeedSpec   `json:"spec,omitempty"`
	Status FeedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FeedList contains a list of Feed
type FeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Feed `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
