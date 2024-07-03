package templates

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
	"time"
)

func TestHighlight(t *testing.T) {
	tests := []struct {
		content  string
		keywords []string
		expected string
	}{
		{"hello world", []string{"world"}, "hello [!]world[!]"},
		{"hello world", []string{"hello", "world"}, "[!]hello[!] [!]world[!]"},
		{"hello world", nil, "hello world"},
		{"", []string{"world"}, ""},
		{"hello", []string{"world"}, "hello"},
	}

	for _, test := range tests {
		result := highlight(test.content, test.keywords)
		if result != test.expected {
			t.Errorf("highlight(%q, %v) = %q; want %q", test.content, test.keywords, result, test.expected)
		}
	}
}

// Test for formatDate function
func TestFormatDate(t *testing.T) {
	tests := []struct {
		time     time.Time
		layout   string
		expected string
	}{
		{time.Date(2023, 6, 18, 12, 0, 0, 0, time.UTC), "2006-01-02", "2023-06-18"},
		{time.Date(2023, 6, 18, 12, 0, 0, 0, time.UTC), "02-Jan-2006", "18-Jun-2023"},
		{time.Date(2023, 6, 18, 12, 0, 0, 0, time.UTC), time.RFC1123, "Sun, 18 Jun 2023 12:00:00 UTC"},
	}

	for _, test := range tests {
		result := formatDate(test.time, test.layout)
		if result != test.expected {
			t.Errorf("formatDate(%v, %q) = %q; want %q", test.time, test.layout, result, test.expected)
		}
	}
}

// Test for contains function
func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		arr      []string
		expected bool
	}{
		{"hello world", []string{"world"}, true},
		{"hello world", []string{"world", "hello"}, true},
		{"hello world", []string{"test"}, false},
		{"hello world", []string{}, false},
		{"", []string{"world"}, false},
	}

	for _, test := range tests {
		result := contains(test.s, test.arr)
		if result != test.expected {
			t.Errorf("contains(%q, %v) = %v; want %v", test.s, test.arr, result, test.expected)
		}
	}
}

func TestPrintTemplate(t *testing.T) {
	testCases := []struct {
		Name  string
		Input struct {
			Filters  *types.FilteringParams
			Articles []types.News
		}
	}{
		{
			Name: "Successful template execution",
			Input: struct {
				Filters  *types.FilteringParams
				Articles []types.News
			}{
				Filters: types.NewFilteringParams("", "", "", ""),
				Articles: []types.News{
					{
						Title:       "Article 1",
						Description: "Description 1",
					},
					{
						Title:       "Article 2",
						Description: "Description 2",
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		err := PrintTemplate(tt.Input.Filters, tt.Input.Articles)

		assert.Nil(t, err)
	}
}
