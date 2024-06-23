package parsers

import (
	"fmt"
	"newsaggr/cmd/types"
	"sync"
	"time"
)

// GenerateDateRange generates an array of date strings between two dates (inclusive).
// The input dates should be in the format "YYYY-MM-DD".
// If the start date is after the end date or if the input date strings are not valid,
// an error is returned.
func GenerateDateRange(dateFrom, dateEnd string) ([]string, error) {
	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, dateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}
	endDate, err := time.Parse(layout, dateEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// Ensure the start date is before or equal to the end date
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	var dates []string
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(layout))
	}

	return dates, nil
}

// FromFiles retrieves news articles from JSON files within the specified date range.
// The date range is inclusive and should be provided in the format "YYYY-MM-DD".
// The function concurrently parses JSON files for each date in the range.
// If an error occurs during the parsing of any file, the process is aborted and the error is returned.
// The returned slice contains all successfully parsed news articles.
func FromFiles(dateFrom, dateEnd string) ([]types.News, error) {
	var (
		news       []types.News
		wg         sync.WaitGroup
		mu         sync.Mutex
		errChannel = make(chan error)
	)

	collectNews := func(p Parser) {
		defer wg.Done()
		n, err := p.Parse()
		if err != nil {
			errChannel <- err
			return
		}
		mu.Lock()
		news = append(news, n...)
		mu.Unlock()
	}

	dates, err := GenerateDateRange(dateFrom, dateEnd)
	if err != nil {
		return nil, err
	}

	for _, date := range dates {
		jp := g.CreateJsonParser(date + ".json")
		wg.Add(1)
		go collectNews(jp)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChannel)

	// Check if there were any errors
	for err := range errChannel {
		if err != nil {
			return nil, err
		}
	}

	return news, nil
}
