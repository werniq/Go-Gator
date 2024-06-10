package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddFetchNewsCmd(t *testing.T) {
	// Verify the command properties
	assert.Equal(t, "fetch", fetchNews.Use, "Command use should be 'fetch'")
	assert.Equal(t, "Fetching news from downloaded data", fetchNews.Short, "Command short description should match")
	assert.Contains(t, fetchNews.Long, "This command parses HTML, XML and JSON files sorts them by given arguments", "Command long description should contain 'This command parses HTML, XML and JSON files'")

	// Verify the flags
	assert.NotNil(t, fetchNews.Flags().Lookup("keywords"), "Flag 'keywords' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("date-from"), "Flag 'date-from' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("date-end"), "Flag 'date-end' should be defined")
	assert.NotNil(t, fetchNews.Flags().Lookup("sources"), "Flag 'sources' should be defined")
}
