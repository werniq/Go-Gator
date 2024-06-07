package templates

import (
	"newsAggr/cmd/filters"
	"newsAggr/cmd/types"
	"sort"
)

type ByPubDate []types.News

// Implement the sort.Interface for ByPubDate
func (a ByPubDate) Len() int      { return len(a) }
func (a ByPubDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPubDate) Less(i, j int) bool {
	//return a[i].PubDate > a[j].PubDate

	t1, err1 := filters.ParseDate(a[i].PubDate)
	t2, err2 := filters.ParseDate(a[j].PubDate)

	// If there is an error parsing dates, consider them equal (or handle errors appropriately)
	if err1 != nil || err2 != nil {
		return false
	}

	return t1.Before(t2)
}

// Function to sort News by PubDate
func sortNewsByPubDate(news []types.News) {
	sort.Sort(ByPubDate(news))
}
