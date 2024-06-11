package parsers

import (
	"newsAggr/cmd/filters"
	"newsAggr/cmd/types"
)

// ApplyParams filters news by provided FilteringParams
func ApplyParams(articles []types.News, params *types.FilteringParams) []types.News {
	var filteredArticles []types.News
	f := filters.InstructionFactory{}

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
