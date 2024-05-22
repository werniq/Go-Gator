package parsers

import (
	"newsAggr/cmd/parsingInstructions"
	"newsAggr/cmd/types"
)

// ApplyParams filters news by provided ParsingParams
func ApplyParams(articles []types.News, params *types.ParsingParams, factory parsingInstructions.InstructionsFactory) []types.News {
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
