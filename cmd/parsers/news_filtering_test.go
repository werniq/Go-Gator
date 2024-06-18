package parsers

import (
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/types"
	"testing"
)

func TestApplyParams(t *testing.T) {
	testCases := []struct {
		Input struct {
			News   []types.News
			Params *types.FilteringParams
		}
		ExpectedOutput []types.News
	}{
		{
			Input: struct {
				News   []types.News
				Params *types.FilteringParams
			}{
				News: []types.News{
					{
						Title:   "Russia's glide bombs devastating Ukraine's cities on the cheap",
						PubDate: "Sun, 19 May 2024 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/cz5drkr8l1ko",
					},
					{
						Title:   "News Title 1",
						PubDate: "Sun, 26 May 2021 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/fake-article-link-1",
					},
					{
						Title:   "News Title 2",
						PubDate: "Sun, 12 May 2021 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/fake-article-link-2",
					},
				},
				Params: &types.FilteringParams{
					Keywords: "glide",
				},
			},
			ExpectedOutput: []types.News{
				{
					Title:   "Russia's glide bombs devastating Ukraine's cities on the cheap",
					PubDate: "Sun, 19 May 2024 07:05:28 GMT",
					Link:    "https://www.bbc.com/news/articles/cz5drkr8l1ko",
				},
			},
		},
	}

	for _, testCase := range testCases {
		var news []types.News
		news = append(news, ApplyFilters(testCase.Input.News, testCase.Input.Params)...)

		assert.Equal(t, news, testCase.ExpectedOutput)
	}
}
