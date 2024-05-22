package cmd

import (
	"fmt"
	"newsAggr/cmd/types"
	"strings"
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

	keywords := strings.Split(params.Keywords, ",")

	var filteredArticles []types.News

	timeFormats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, article := range articles {
		keywordOk := false
		for _, keyword := range keywords {
			if strings.Contains(article.Title, keyword) || strings.Contains(article.Description, keyword) {
				keywordOk = true
				break
			}
		}
		if !keywordOk {
			continue
		}

		var publicationDate time.Time
		var err error

		if article.PubDate != "" {
			publicationDate, err = parseDateWithFormats(article.PubDate, timeFormats)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if params.StartingTimestamp != "" {
			startingTime, err := parseDateWithFormats(params.StartingTimestamp, timeFormats)
			if err != nil {
				continue
			}

			if publicationDate.Before(startingTime) {
				continue
			}
		}

		if params.EndingTimestamp != "" {
			endingTime, err := parseDateWithFormats(params.EndingTimestamp, timeFormats)
			if err != nil {
				continue
			}

			if publicationDate.After(endingTime) {
				continue
			}
		}

		filteredArticles = append(filteredArticles, article)
	}

	return filteredArticles
}
