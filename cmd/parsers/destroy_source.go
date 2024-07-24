package parsers

import (
	"encoding/json"
	"gogator/cmd/types"
	"os"
	"path/filepath"
)

// DestroySource will be called whenever we delete source from database
//
// This function removes articles with this source as publisher, from all available data files
func DestroySource(source string, dateRange []string) error {
	for _, filename := range dateRange {
		jp := JsonParser{
			Source: filename + JsonExtension,
		}

		n, err := jp.Parse()
		if err != nil {
			return nil
		}

		n = removeNewsBySource(n, source)
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		filename = filepath.Join(cwd, StoragePath, filename+JsonExtension)

		err = os.Truncate(filename, 0)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(filename, os.O_RDWR, 0654)
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
func removeNewsBySource(n []types.News, source string) []types.News {
	var news []types.News
	for _, article := range n {
		if article.Publisher != source {
			news = append(news, article)
		}
	}

	return news
}
