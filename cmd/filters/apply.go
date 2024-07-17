package filters

import "gogator/cmd/types"

var (
	f InstructionFactory
)

// Apply filters news by provided FilteringParams
func Apply(articles []types.News, params *types.FilteringParams) []types.News {
	var filteredArticles []types.News

	filters := []func(article types.News, params *types.FilteringParams) bool{
		f.CreateSourcesInstruction().Apply,
		f.CreateApplyDataRangeInstruction().Apply,
		f.CreateApplyKeywordInstruction().Apply,
	}

	// Iterate over the articles
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
