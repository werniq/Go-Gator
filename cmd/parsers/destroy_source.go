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
	for _, articlesFilename := range dateRange {
		jp := JsonParser{
			Source: articlesFilename + JsonExtension,
		}

		news, err := jp.Parse()
		if err != nil {
			return nil
		}

		news = removeNewsBySource(news, source)

		cwdPath, err := os.Getwd()
		if err != nil {
			return err
		}

		articlesFilename = filepath.Join(cwdPath, StoragePath, articlesFilename+JsonExtension)

		err = os.Truncate(articlesFilename, 0)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(articlesFilename, os.O_RDWR, 0654)
		if err != nil {
			return err
		}

		articleFileData, err := json.Marshal(news)
		if err != nil {
			return err
		}

		_, err = file.Write(articleFileData)
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
func removeNewsBySource(news []types.Article, source string) []types.Article {
	var filteredNews []types.Article
	for _, article := range news {
		if article.Publisher != source {
			filteredNews = append(filteredNews, article)
		}
	}

	return filteredNews
}
