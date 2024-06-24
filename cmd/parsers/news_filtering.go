package parsers

import (
	"newsaggr/cmd/types"
	"strings"
)

var PublisherMapping = map[string]string{
	ABC:             ABC,
	BBC:             BBC,
	WashingtonTimes: WashingtonTimes,
	UsaToday:        UsaToday,
}

// determinePublisher returns the publisher based on the URL
func determinePublisher(url string) string {
	for key, publisher := range PublisherMapping {
		if strings.Contains(url, key) {
			return publisher
		}
	}
	return "default"
}

// ApplyFilters filters news by provided FilteringParams
func ApplyFilters(articles []types.News, params *types.FilteringParams) []types.News {
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
