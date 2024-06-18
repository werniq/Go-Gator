package parsers

import (
	"newsaggr/cmd/types"
)

// ApplyFilters filters news by provided FilteringParams
func ApplyFilters(articles []types.News, params *types.FilteringParams) []types.News {
	var filteredArticles []types.News

	keywordInstruction := f.CreateApplyKeywordInstruction()
	dateRangeInstruction := f.CreateApplyDataRangeInstruction()

	for _, article := range articles {
		if !keywordInstruction.Apply(article, params) {
			continue
		}
		if !dateRangeInstruction.Apply(article, params) {
			continue
		}
		filteredArticles = append(filteredArticles, article)
	}

	return filteredArticles
}
