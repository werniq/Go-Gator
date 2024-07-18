package parsers

import "gogator/cmd/types"

// DestroySource will be called whenever we delete source from database
// All articles with this source as publisher will be deleted
func DestroySource(source string, dateRange []string) error {
	for _, filename := range dateRange {
		jp := JsonParser{
			Source: filename + ".json",
		}

		n, err := jp.Parse()
		if err != nil {
			return nil
		}

		n, err = removeNewsBySource(n, source)
		if err != nil {
			return err
		}
	}

	return nil
}

// removeNewsBySource removes all news items with the specified source
func removeNewsBySource(n []types.News, source string) ([]types.News, error) {
	for i := 0; i <= len(n)-1; i++ {
		if n[i].Publisher == source {
			n = append(n[:i], n[i+1:]...)
		}
	}

	return n, nil
}
