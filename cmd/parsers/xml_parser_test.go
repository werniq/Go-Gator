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
			Expected: []types.News{},
			Input:    &types.FilteringParams{},
		},
	}

	for _, testCase := range testCases {
		news, err := parser.Parse()
		assert.NoError(t, err)

		filteredNews := ApplyFilters(news, testCase.Input)

		assert.NotNil(t, filteredNews)
	}
}
