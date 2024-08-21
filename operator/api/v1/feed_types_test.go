package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetCreatedCondition(t *testing.T) {
	feed := &Feed{}
	feed.SetCreatedCondition("Initial setup")

	condition, exists := feed.Status.Conditions[TypeFeedCreated]
	assert.True(t, exists)
	assert.Equal(t, createdReason, condition.Status)
	assert.Equal(t, "Initial setup", condition.Reason)
	assert.Equal(t, FeedCreated, condition.Message)
}

func TestSetFailedCondition(t *testing.T) {
	feed := &Feed{}
	feed.SetFailedCondition("Resource quota exceeded", "QuotaExceeded")

	condition, exists := feed.Status.Conditions[TypeFeedFailedToCreate]
	assert.True(t, exists)
	assert.Equal(t, failedToCreateReason, condition.Status)
	assert.Equal(t, "QuotaExceeded", condition.Reason)
	assert.Equal(t, "Resource quota exceeded", condition.Message)
}

func TestSetUpdatedCondition(t *testing.T) {
	feed := &Feed{}
	feed.SetUpdatedCondition("Updated after scaling")

	condition, exists := feed.Status.Conditions[TypeFeedUpdated]
	assert.True(t, exists)
	assert.Equal(t, createdReason, condition.Status)
	assert.Equal(t, "Updated after scaling", condition.Reason)
	assert.Equal(t, FeedUpdated, condition.Message)
}
