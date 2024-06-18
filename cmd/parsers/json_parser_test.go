package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestJsonParser_ParseWithArgs(t *testing.T) {
	parser := JsonParser{
		Source: "nbc",
	}

	testCases := []struct {
		Expected []types.News
		Input    *types.FilteringParams
	}{
		{
			Expected: []types.News{
				{
					Title:       "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena",
					Description: "The Vatican has issued new rules radically reforming its process for evaluating faith-based supernatural phenomena like visions of the Virgin Mary or stigmata.",
					PubDate:     "2024-05-17T14:58:52Z",
				},
			},
			Input: &types.FilteringParams{
				Keywords: "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena",
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
				assert.Equal(t, filteredNews[i].PubDate, expectedNews.PubDate)
			}
		}
	}
}
