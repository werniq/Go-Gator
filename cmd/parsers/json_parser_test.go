package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/types"
	"testing"
)

func TestJsonParser_ParseWithArgs(t *testing.T) {
	parser := JsonParser{
		Source: "2024-06-20.json",
	}

	testCases := []struct {
		Expected []types.News
		Input    *types.FilteringParams
	}{
		{
			Expected: []types.News{
				{
					Title:       "Historic flooding in southern China kills 47, with more floods feared in coming days\n",
					Description: "At least 47 people have died as downpours in southern China's Guangdong province caused historic flooding and slides, state media reported Friday, while authorities warned of more extreme weather ahead in other parts of the country.\n",
					PubDate:     "Fri, 21 Jun 2024 07:26:14 -0400",
				},
			},
			Input: &types.FilteringParams{
				Keywords: "Historic flooding",
			},
		},
	}

	for _, testCase := range testCases {
		news, err := parser.Parse()
		assert.NoError(t, err)

		filteredNews := filters.Apply(news, testCase.Input)

		if len(testCase.Expected) == 0 {
			assert.Empty(t, filteredNews)
		} else {
			assert.Equal(t, len(filteredNews), len(testCase.Expected))
			for i, expectedNews := range testCase.Expected {
				assert.Equal(t, filteredNews[i].Title, expectedNews.Title)
				assert.Equal(t, filteredNews[i].Description, expectedNews.Description)
				assert.Equal(t, filteredNews[i].PubDate, expectedNews.PubDate)
			}
		}
	}
}
