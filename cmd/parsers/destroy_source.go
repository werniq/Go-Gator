package parsers

import (
	"encoding/json"
	"gogator/cmd/types"
	"os"
)

// DestroySource will be called whenever we delete source from database
// All articles with this source as publisher will be deleted
func DestroySource(source string, dateRange []string) error {
	for _, filename := range dateRange {
		jp := JsonParser{
			Source: filename + JsonExtension,
		}

		n, err := jp.Parse()
		if err != nil {
			return nil
		}

		n, err = removeNewsBySource(n, source)
		if err != nil {
			return err
		}

		file, err := os.Open(filename + JsonExtension)
		if err != nil {
			return err
		}

		out, err := json.Marshal(n)
		if err != nil {
			return err
		}

		_, err = file.Write(out)
		if err != nil {
			return err
		}

		err = file.Close()
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
