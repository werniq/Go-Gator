package controller

import (
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/event"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	"testing"
)

func TestFeedStatusConditionPredicate_Update(t *testing.T) {
	predicate := FeedStatusConditionPredicate{}

	tests := []struct {
		name     string
		e        event.UpdateEvent
		expected bool
	}{
		{
			name:     "ObjectNew is nil",
			e:        event.UpdateEvent{ObjectNew: nil},
			expected: false,
		},
		{
			name: "ObjectNew is *newsaggregatorv1.Feed with Conditions['Created'] true",
			e: event.UpdateEvent{ObjectNew: &newsaggregatorv1.Feed{
				Status: newsaggregatorv1.FeedStatus{
					Conditions: map[string]newsaggregatorv1.FeedConditions{
						"Created": {Status: true},
					},
				},
			}},
			expected: true,
		},
		{
			name: "ObjectNew is *newsaggregatorv1.Feed with Conditions['Deleted'] true",
			e: event.UpdateEvent{ObjectNew: &newsaggregatorv1.Feed{
				Status: newsaggregatorv1.FeedStatus{
					Conditions: map[string]newsaggregatorv1.FeedConditions{
						"Deleted": {Status: true},
					},
				},
			}},
			expected: true,
		},
		{
			name: "ObjectNew is *newsaggregatorv1.Feed but neither Conditions['Created'] nor Conditions['Deleted'] are true",
			e: event.UpdateEvent{ObjectNew: &newsaggregatorv1.Feed{
				Status: newsaggregatorv1.FeedStatus{
					Conditions: map[string]newsaggregatorv1.FeedConditions{
						"Created": {Status: false},
						"Deleted": {Status: false},
					},
				},
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicate.Update(tt.e)
			assert.Equal(t, tt.expected, result)
		})
	}
}
