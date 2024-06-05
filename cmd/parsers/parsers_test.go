package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsAggr/cmd/types"
	"testing"
)

func TestParseWithParams(t *testing.T) {
	formats := []string{"xml", "json", "html"}
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
			ExpectedOutput: 9,
		},
	}

	for _, testCase := range testCases {
		var news []types.News
		for _, format := range formats {
			news = append(news, ParseWithParams(format, testCase.Input)...)
		}
		assert.Equal(t, testCase.ExpectedOutput, len(news))
	}
}
