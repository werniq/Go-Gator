package filters

import (
	"newsAggr/cmd/types"
	"strings"
	"time"
)

type ApplyKeywordsInstruction struct{}

// Apply is a method in ApplyKeywordsInstruction which is used to filter articles by keyword
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

// Apply in ApplyDateRangeInstruction is a method which is used to filter articles by data range
func (a ApplyDateRangeInstruction) Apply(article types.News, params *types.FilteringParams) bool {
	timeFormats := []string{
		time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z,
		time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
		time.DateTime, time.DateOnly, time.TimeOnly, "January 2, 2006",
	}

	var publicationDate time.Time
	var err error

	if article.PubDate != "" {
		publicationDate, err = parseDate(article.PubDate, timeFormats)
		if err != nil {
			return false
		}
	}

	if params.StartingTimestamp != "" {
		startingTime, err := parseDate(params.StartingTimestamp, timeFormats)
		if err != nil {
			return false
		}
		if publicationDate.Before(startingTime) {
			return false
		}
	}

	if params.EndingTimestamp != "" {
		endingTime, err := parseDate(params.EndingTimestamp, timeFormats)
		if err != nil {
			return false
		}
		if publicationDate.After(endingTime) {
			return false
		}
	}

	return true
}

// parseDate is utility function which parses date to time.Time object
func parseDate(dateStr string, formats []string) (time.Time, error) {
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
