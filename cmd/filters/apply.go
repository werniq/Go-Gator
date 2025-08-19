package filters

import "gogator/cmd/types"

var (
	f InstructionFactory
)

// Apply filters news by provided FilteringParams
func Apply(articles []types.Article, params *types.FilteringParams) []types.Article {
	var filteredArticles []types.Article

	filters := []func(article types.Article, params *types.FilteringParams) bool{
		f.CreateSourcesInstruction().Apply,
		f.CreateApplyDataRangeInstruction().Apply,
		f.CreateApplyKeywordInstruction().Apply,
	}

	for _, article := range articles {
		applyAllFilters := true

		for _, filter := range filters {
			if !filter(article, params) {
				applyAllFilters = false
				break
			}
		}
		if applyAllFilters {
			filteredArticles = append(filteredArticles, article)
		}
	}

	return filteredArticles
}
