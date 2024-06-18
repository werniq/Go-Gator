package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestXmlParser_Parse(t *testing.T) {
	parser := XMLParser{
		Source: "abc",
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
				Keywords: "Ukraine and Russia exchange drone attacks while Russia continues its push in the east",
			},
		},
		{
			Expected: []types.News{
				{
					Title:       "Statue weeping blood? Visions of the Virgin Mary? Vatican has new advice on supernatural phenomena",
					Description: "The Vatican has issued new rules radically reforming its process for evaluating faith-based supernatural phenomena like visions of the Virgin Mary or stigmata.",
					PubDate:     "2024-05-17T14:58:52Z",
				},
			},
			Input: &types.FilteringParams{
				Keywords: "Definitely not existen sequence of keywords",
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
			assert.Equal(t, len(testCase.Expected), len(filteredNews))
		}
	}
}
