package types

// FilteringParams represents the parameters used for filtering news articles.
// It has several fields:
// /  1. Keywords          - Keywords to filter articles
// /  2. StartingTimestamp - Starting timestamp for filtering articles
// /  3. EndingTimestamp   - Ending timestamp for filtering articles
// /  4. Sources           - Sources to filter articles
//
// This struct will be used for:
//  1. Handling user input
//  2. Filter news by these criteria
type FilteringParams struct {
	Keywords          string `json:"keywords" xml:"keywords"`
	StartingTimestamp string `json:"starting_timestamp" xml:"starting_timestamp"`
	EndingTimestamp   string `json:"ending_timestamp" xml:"ending_timestamp"`
	Sources           string `json:"sources" xml:"sources"`
}

// NewFilteringParams creates an instance of FilteringParams
func NewFilteringParams(keywords, start, end, sources string) *FilteringParams {
	return &FilteringParams{
		Keywords:          keywords,
		StartingTimestamp: start,
		EndingTimestamp:   end,
		Sources:           sources,
	}
}
