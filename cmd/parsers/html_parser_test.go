package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/types"
	"testing"
)

func TestHtmlParser_Parse(t *testing.T) {
	parser := HtmlParser{
		Source: UsaToday,
	}

	testCases := []struct {
		Expected []types.News
		Input    *types.FilteringParams
	}{
		{
			Expected: []types.News{},
			Input:    &types.FilteringParams{},
		},
	}

	for _, testCase := range testCases {
		news, err := parser.Parse()
		assert.NoError(t, err)

		filteredNews := filters.Apply(news, testCase.Input)

		assert.NotNil(t, filteredNews)
	}
}
