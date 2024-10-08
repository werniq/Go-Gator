package filters

import (
	"gogator/cmd/types"
	"strings"
	"time"
)

type ApplyKeywordsInstruction struct{}

// Apply is a method in ApplyKeywordsInstruction which is used to filter article by keyword
func (a ApplyKeywordsInstruction) Apply(article types.News, params *types.FilteringParams) bool {
	keywords := strings.Split(params.Keywords, ",")
	for _, keyword := range keywords {
		if strings.Contains(article.Title, keyword) || strings.Contains(article.Description, keyword) {
			return true
		}
	}

	return false
}

type ApplyDateRangeInstruction struct{}

// Apply in ApplyDateRangeInstruction is a method which is used to filter article by data range
func (a ApplyDateRangeInstruction) Apply(article types.News, params *types.FilteringParams) bool {

	var publicationDate time.Time
	var err error

	if article.PubDate != "" {
		publicationDate, err = parseDate(article.PubDate)
		if err != nil {
			return false
		}
	}

	if params.StartingTimestamp != "" {
		startingTime, err := parseDate(params.StartingTimestamp)
		if err != nil {
			return false
		}
		if publicationDate.Before(startingTime) {
			return false
		}
	}

	if params.EndingTimestamp != "" {
		endingTime, err := parseDate(params.EndingTimestamp)
		if err != nil {
			return false
		}
		if publicationDate.After(endingTime) {
			return false
		}
	}

	return true
}

type ApplySourcesInstruction struct{}

// Apply in ApplySourcesInstruction is a method which is used to filter articles by publisher
func (a ApplySourcesInstruction) Apply(article types.News, params *types.FilteringParams) bool {
	if params.Sources == "" {
		return true
	}

	sources := strings.Split(params.Sources, ",")
	for _, source := range sources {
		if article.Publisher == source {
			return true
		}
	}

	return false
}

// parseDate is utility function which parses date to time.Time object
func parseDate(dateStr string) (time.Time, error) {
	var err error
	var date time.Time
	timeFormats := []string{
		time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z,
		time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
		time.DateTime, time.DateOnly, time.TimeOnly, "January 2, 2006",
	}

	for _, format := range timeFormats {
		date, err = time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}
	return time.Time{}, err
}
