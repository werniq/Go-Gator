package types

type FilteringParams struct {
	Keywords          string   `json:"keywords" xml:"keywords"`
	StartingTimestamp string   `json:"starting_timestamp" xml:"starting_timestamp"`
	EndingTimestamp   string   `json:"ending_timestamp" xml:"ending_timestamp"`
	Sources           []string `json:"sources" xml:"sources"`
}
