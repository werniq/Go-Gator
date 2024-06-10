package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
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
			ExpectedOutput: 1,
		},
		{
			Input: &types.FilteringParams{
				Keywords: "Ukraine",
			},
			ExpectedOutput: 3,
		},
	}

	for _, testCase := range testCases {
		var news []types.News
		for _, source := range sources {
			news = append(news, ApplyParams(ParseBySource(source), testCase.Input)...)
		}
		assert.Equal(t, testCase.ExpectedOutput, len(news))
	}
}
