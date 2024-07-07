package templates

import (
	"newsaggr/cmd/types"
	"sort"
	"time"
)

// ByPubDate is a type alias for a slice of types.News, used for sorting purposes
type ByPubDate []types.News

// Len is part of the sort.Interface implementation for ByPubDate
// It returns the number of elements in the slice
func (a ByPubDate) Len() int {
	return len(a)
}

// Swap is part of the sort.Interface implementation for ByPubDate
// It swaps the elements with indexes i and j
func (a ByPubDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less is part of the sort.Interface implementation for ByPubDate
// It returns true if the PubDate of the element at index i is before the PubDate of the element at index j
func (a ByPubDate) Less(i, j int) bool {
	// Parse the publication dates of the elements at indexes i and j.
	t1, err1 := parseDate(a[i].PubDate)
	t2, err2 := parseDate(a[j].PubDate)

	// If there is an error parsing either date, consider them equal
	// (alternatively, you could handle errors in a different way).
	if err1 != nil || err2 != nil {
		return false
	}

	// Return true if t1 is before t2.
	return t1.Before(t2)
}

// sortNewsByPubDate sorts a slice of types.News by their PubDate in ascending order
func sortNewsByPubDate(news []types.News) {
	sort.Sort(ByPubDate(news))
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
