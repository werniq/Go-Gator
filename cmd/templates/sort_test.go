package templates

import (
	"newsaggr/cmd/types"
	"testing"
)

func TestByPubDate_Swap(t *testing.T) {
	news := ByPubDate{
		{PubDate: "2023-05-01"},
		{PubDate: "2022-04-01"},
	}

	news.Swap(0, 1)

	expected := ByPubDate{
		{PubDate: "2022-04-01"},
		{PubDate: "2023-05-01"},
	}

	for i := range news {
		if news[i].PubDate != expected[i].PubDate {
			t.Errorf("Swap() = %v, expected %v", news, expected)
			break
		}
	}
}

// Test the Less method
func TestByPubDate_Less(t *testing.T) {
	news := ByPubDate{
		{PubDate: "2023-05-01"},
		{PubDate: "2022-04-01"},
	}

	if got := news.Less(0, 1); got != false {
		t.Errorf("Less() = %v, expected false", got)
	}

	if got := news.Less(1, 0); got != true {
		t.Errorf("Less() = %v, expected true", got)
	}
}

// Test the Less method with parsing errors
func TestByPubDate_Less_WithParsingErrors(t *testing.T) {
	news := ByPubDate{
		{PubDate: "2024-05-12"},
		{PubDate: "2022-04-01"},
	}

	if got := news.Less(0, 1); got != false {
		t.Errorf("Less() with parsing error = %v, expected %v", got, false)
	}

	if got := news.Less(1, 0); got != true {
		t.Errorf("Less() with parsing error = %v, expected %v", got, true)
	}
}

func TestSortNewsByPubDate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		news     []types.News
		expected []types.News
	}{
		{
			name: "Simple case",
			news: []types.News{
				{PubDate: "2023-05-01"},
				{PubDate: "2022-04-01"},
				{PubDate: "2023-01-01"},
			},
			expected: []types.News{
				{PubDate: "2022-04-01"},
				{PubDate: "2023-01-01"},
				{PubDate: "2023-05-01"},
			},
		},
		{
			name:     "Empty slice",
			news:     []types.News{},
			expected: []types.News{},
		},
		{
			name: "Already sorted",
			news: []types.News{
				{PubDate: "2022-01-01"},
				{PubDate: "2022-02-01"},
				{PubDate: "2022-03-01"},
			},
			expected: []types.News{
				{PubDate: "2022-01-01"},
				{PubDate: "2022-02-01"},
				{PubDate: "2022-03-01"},
			},
		},
		{
			name: "With parsing errors",
			news: []types.News{
				{PubDate: "invalid-date"},
				{PubDate: "2022-01-01"},
			},
			expected: []types.News{
				{PubDate: "invalid-date"},
				{PubDate: "2022-01-01"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortNewsByPubDate(tt.news)
			for i := range tt.news {
				if tt.news[i].PubDate != tt.expected[i].PubDate {
					t.Errorf("expected %v, got %v", tt.expected, tt.news)
					break
				}
			}
		})
	}
}
