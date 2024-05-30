package parsers

import (
	"newsAggr/cmd/filters"
	"newsAggr/cmd/types"
)

// ApplyParams filters news by provided FilteringParams
func ApplyParams(articles []types.News, params *types.FilteringParams, factory filters.InstructionsFactory) []types.News {
	var filteredArticles []types.News

	keywordInstruction := factory.CreateApplyKeywordInstruction()
	dateRangeInstruction := factory.CreateApplyDataRangeInstruction()

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
