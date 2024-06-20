package parsers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestParseWithParams(t *testing.T) {
	sources := []string{"abc", "bbc"}
	testCases := []struct {
		Input          *types.FilteringParams
		ExpectedOutput int
	}{
		{
			Input: &types.FilteringParams{
				Keywords: "glide",
			},
			ExpectedOutput: 0,
		},
		{
			Input: &types.FilteringParams{
				Keywords: "Ukraine",
			},
			ExpectedOutput: 2,
		},
	}

	for _, testCase := range testCases {
		var news []types.News
		var err error

		for _, source := range sources {
			news, err = ParseBySource(source)
			assert.Equal(t, err, nil, fmt.Sprintf("Expected: %v, Got: %v", nil, err))
		}
		news = ApplyFilters(news, testCase.Input)
		assert.Equal(t, testCase.ExpectedOutput, len(news))
	}
}
