package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		Source: "usatoday",
	}

	testCases := []struct {
		Expected []types.News
		Input    *types.FilteringParams
	}{
		{
			Expected: []types.News{},
			Input: &types.FilteringParams{
				Keywords: "Mac address",
			},
		},
		{
			Expected: []types.News{
				{
					Title:       "The Taizhou Zoo in Jiangsu, China dyed two chow chow dogs and advertised them as \"panda dogs\" in the exhibit that opened on May 1..",
					Description: "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit",
				},
			},
			Input: &types.FilteringParams{
				Keywords: "Zoo in China criticized for dyeing dogs for 'panda dog' exhibit",
			},
		},
		{
			Expected: []types.News{},
			Input: &types.FilteringParams{
				Keywords: "definitely not-existent sequence of keywords",
			},
		},
	}

	for _, testCase := range testCases {
		news, err := parser.Parse()
		assert.NoError(t, err)

		filteredNews := ApplyFilters(news, testCase.Input)

		if len(testCase.Expected) == 0 {
			assert.Empty(t, filteredNews)
		} else {
			assert.Equal(t, len(filteredNews), len(testCase.Expected))
			for i, expectedNews := range testCase.Expected {
				assert.Equal(t, filteredNews[i].Title, expectedNews.Title)
				assert.Equal(t, filteredNews[i].Description, expectedNews.Description)
			}
		}
	}
}
