package parsers

import (
	"encoding/json"
	"fmt"
	"newsAggr/cmd/FilteringInstructions"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"reflect"
)

type JsonParser struct {
}

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse(params *types.FilteringParams) []types.News {
	var news []types.News

	sourceToFile := map[string]string{
		"nbc": "nbc-news.json",
	}

	var filenames []string

	if !reflect.DeepEqual(params.Sources, filenames) {
		for _, val := range params.Sources {
			if filename, ok := sourceToFile[val]; ok {
				fmt.Println(filename)
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
		data := ExtractFileData(filename)
		if data == nil {
			logger.ErrorLogger.Fatalf("Error extracting file data: %v\n", filename)
		}

		var dummy types.Json
		err := json.Unmarshal(data, &dummy)
		if err != nil {
			logger.ErrorLogger.Fatalf("Error decoding JSON data: %v\n", err)
		}

		news = append(news, types.JsonNewsToNews(dummy.Articles)...)
	}

	factory := FilteringInstructions.GoGatorInstructionFactory{}

	return ApplyParams(news, params, factory)
}
