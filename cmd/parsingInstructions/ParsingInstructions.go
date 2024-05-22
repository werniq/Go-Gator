package parsingInstructions

import (
	"newsAggr/cmd/types"
	"newsAggr/cmd/utils"
	"strings"
	"time"
)

type ApplyKeywordsInstruction struct{}

func (a ApplyKeywordsInstruction) Apply(article types.News, params *types.ParsingParams) bool {
	keywords := strings.Split(params.Keywords, ",")
	for _, keyword := range keywords {
		if strings.Contains(article.Title, keyword) || strings.Contains(article.Description, keyword) {
			return true
		}
	}
	return false
}

type ApplyDateRangeInstruction struct{}

func (a ApplyDateRangeInstruction) Apply(article types.News, params *types.ParsingParams) bool {
	timeFormats := []string{
		time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z,
		time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
		time.DateTime, time.DateOnly, time.TimeOnly, "January 2, 2006",
	}

	var publicationDate time.Time
	var err error

	if article.PubDate != "" {
		publicationDate, err = utils.ParseDateWithFormats(article.PubDate, timeFormats)
		if err != nil {
			return false
		}
	}

	if params.StartingTimestamp != "" {
		startingTime, err := utils.ParseDateWithFormats(params.StartingTimestamp, timeFormats)
		if err != nil {
			return false
		}
		if publicationDate.Before(startingTime) {
			return false
		}
	}

	if params.EndingTimestamp != "" {
		endingTime, err := utils.ParseDateWithFormats(params.EndingTimestamp, timeFormats)
		if err != nil {
			return false
		}
		if publicationDate.After(endingTime) {
			return false
		}
	}

	return true
}
