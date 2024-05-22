package cmd

import (
	"newsAggr/cmd/types"
	"time"
)

type ParsingParams struct {
	Keywords          string   `json:"keywords" xml:"keywords"`
	StartingTimestamp string   `json:"starting_timestamp" xml:"starting_timestamp"`
	EndingTimestamp   string   `json:"ending_timestamp" xml:"ending_timestamp"`
	Sources           []string `json:"sources" xml:"sources"`
}

type Parsers interface {
	Parse(params ParsingParams) []types.News
}

func parseDateWithFormats(dateStr string, formats []string) (time.Time, error) {
	var err error
	var date time.Time

	for _, format := range formats {
		date, err = time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}
	return time.Time{}, err
}

// ApplyParams filters news by provided ParsingParams
func ApplyParams(articles []types.News, params ParsingParams, factory ParsingInstructionsFactory) []types.News {
	if articles == nil {
		return nil
	}

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
