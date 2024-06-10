package parsers

import (
	"encoding/json"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"reflect"
	"sync"
)

type JsonParser struct {
}

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse(params *types.FilteringParams) []types.News {
	var news []types.News
	var wg sync.WaitGroup

	sourceToFile := map[string]string{
		"nbc": "nbc-news.json",
	}

	var filenames []string

	if !reflect.DeepEqual(params.Sources, filenames) {
		for _, val := range params.Sources {
			if filename, ok := sourceToFile[val]; ok {
				filenames = append(filenames, filename)
			}
		}

		if len(filenames) == 0 {
			return nil
		}
	} else {
		filenames = []string{"nbc-news.json"}
	}

	for _, filename := range filenames {
		wg.Add(1)

		go func(filename string) {
			defer wg.Done()
			data := extractFileData(filename)
			if data == nil {
				logger.ErrorLogger.Fatalf("Error extracting file data: %v\n", filename)
			}

			var dummy types.Json
			err := json.Unmarshal(data, &dummy)
			if err != nil {
				logger.ErrorLogger.Fatalf("Error decoding JSON data: %v\n", err)
			}

			news = append(news, types.JsonNewsToNews(dummy.Articles, filename)...)
		}(filename)
	}

	wg.Wait()

	return news
}
