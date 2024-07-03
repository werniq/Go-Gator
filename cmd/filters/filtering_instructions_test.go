package filters

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestApplyKeywordsInstruction_Apply(t *testing.T) {
	testCases := []struct {
		Input struct {
			Article types.News
			Params  *types.FilteringParams
		}
		ExpectedOutput bool
	}{
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					Title: "Ukraine is a great country",
				},
				Params: &types.FilteringParams{
					Keywords: "Ukraine",
				},
			},
			ExpectedOutput: true,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					Title: "Ukraine is a great country",
				},
				Params: &types.FilteringParams{
					Keywords: "China",
				},
			},
			ExpectedOutput: false,
		},
	}
	keywordInstruction := ApplyKeywordsInstruction{}

	for _, testCase := range testCases {
		match := keywordInstruction.Apply(
			testCase.Input.Article,
			testCase.Input.Params)

		assert.Equal(t, match, testCase.ExpectedOutput)
	}
}

func TestApplyDateRangeInstruction_Apply(t *testing.T) {
	testCases := []struct {
		Input struct {
			Article types.News
			Params  *types.FilteringParams
		}
		ExpectedOutput bool
	}{
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-12",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-11",
				},
			},
			ExpectedOutput: true,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-55",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-12",
				},
			},
			ExpectedOutput: false,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-12",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-55",
				},
			},
			ExpectedOutput: false,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-12",
				},
				Params: &types.FilteringParams{
					EndingTimestamp: "2024-05-55",
				},
			},
			ExpectedOutput: false,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-10",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-12",
				},
			},
			ExpectedOutput: false,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-14",
				},
				Params: &types.FilteringParams{
					EndingTimestamp: "2024-05-13",
				},
			},
			ExpectedOutput: false,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-14",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-12",
					EndingTimestamp:   "2024-05-15",
				},
			},
			ExpectedOutput: true,
		},
		{
			Input: struct {
				Article types.News
				Params  *types.FilteringParams
			}{
				Article: types.News{
					PubDate: "2024-05-19",
				},
				Params: &types.FilteringParams{
					StartingTimestamp: "2024-05-12",
					EndingTimestamp:   "2024-05-15",
				},
			},
			ExpectedOutput: false,
		},
	}
	dateRangeInstruction := ApplyDateRangeInstruction{}

	for _, testCase := range testCases {
		match := dateRangeInstruction.Apply(
			testCase.Input.Article,
			testCase.Input.Params)

		assert.Equal(t, match, testCase.ExpectedOutput)
	}
}
