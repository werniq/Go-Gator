package filters

import (
	"github.com/stretchr/testify/assert"
	"gogator/cmd/types"
	"testing"
)

func TestApply(t *testing.T) {
	testCases := []struct {
		Input struct {
			News   []types.Article
			Params *types.FilteringParams
		}
		ExpectedOutput []types.Article
	}{
		{
			Input: struct {
				News   []types.Article
				Params *types.FilteringParams
			}{
				News: []types.Article{
					{
						Title:   "Russia's glide bombs devastating Ukraine's cities on the cheap",
						PubDate: "Sun, 19 May 2024 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/cz5drkr8l1ko",
					},
					{
						Title:   "Article Title 1",
						PubDate: "Sun, 26 May 2021 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/fake-article-link-1",
					},
					{
						Title:   "Article Title 2",
						PubDate: "Sun, 12 May 2021 07:05:28 GMT",
						Link:    "https://www.bbc.com/news/articles/fake-article-link-2",
					},
				},
				Params: &types.FilteringParams{
					Keywords: "glide",
				},
			},
			ExpectedOutput: []types.Article{
				{
					Title:   "Russia's glide bombs devastating Ukraine's cities on the cheap",
					PubDate: "Sun, 19 May 2024 07:05:28 GMT",
					Link:    "https://www.bbc.com/news/articles/cz5drkr8l1ko",
				},
			},
		},
	}

	for _, testCase := range testCases {
		var news []types.Article
		news = append(news, Apply(testCase.Input.News, testCase.Input.Params)...)

		assert.Equal(t, news, testCase.ExpectedOutput)
	}
}
