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

const (
	// TypeFeedCreated represents the Created condition type
	TypeFeedCreated = "Created"

	// TypeFeedFailedToCreate represents the FailedToCreate condition type
	TypeFeedFailedToCreate = "FailedToCreate"

	// TypeFeedUpdated represents the Updated condition type
	TypeFeedUpdated = "Updated"

	// failedToCreateReason is a constant that represents the reason for failed to create condition
	failedToCreateReason = false

	// createdReason is a constant that represents the reason for created condition
	createdReason = true

	// feedStatusConditionsCapacity is a capacity of feed status conditions array
	feedStatusConditionsCapacity = 3

	// FeedCreated represents the reason for created condition
	FeedCreated = "FeedCreated"

	// FeedUpdated represents the reason for updated condition
	FeedUpdated = "FeedUpdated"
)

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	// Name field is a string that represents the name of the feed
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Link field is a string that represents the URL of the feed
	// +kubebuilder:validation:Required
	Link string `json:"link,omitempty"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	// Conditions field is a map of conditions that the feed can have
	// question: i tried to use enum with int iota as mappings keys but i had an error
	// saying that map keys should be strings, not int
	Conditions map[string]FeedConditions `json:"conditions,omitempty"`
}

// FeedConditions are the cond
type FeedConditions struct {
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

// SetCreatedCondition sets the created condition of the feed to Created
func (r *Feed) SetCreatedCondition(reason string) {
	r.setCondition(TypeFeedCreated, createdReason, reason, FeedCreated)
}

// SetFailedCondition sets the failed condition of the feed to Failed
//
// It sets the status to false, the reason to failedToCreateReason and the message with reason to the provided message.
func (r *Feed) SetFailedCondition(message, reason string) {
	r.setCondition(TypeFeedFailedToCreate, failedToCreateReason, reason, message)
}

// SetUpdatedCondition sets the updated condition of the feed to Updated
func (r *Feed) SetUpdatedCondition(reason string) {
	r.setCondition(TypeFeedUpdated, createdReason, reason, FeedUpdated)
}

// setCondition sets the created condition of the feed to the one specified in arguments.
// Is used for aliasing the created and updated conditions.
func (r *Feed) setCondition(conditionType string, status bool, reason, message string) {
	if r.Status.Conditions == nil {
		r.Status.Conditions = make(map[string]FeedConditions, feedStatusConditionsCapacity)
	}

	r.Status.Conditions[conditionType] = FeedConditions{
		Status:         status,
		Reason:         reason,
		Message:        message,
		LastUpdateTime: metav1.Now().String(),
	}
}
