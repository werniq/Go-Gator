package parsers

import (
	"encoding/json"
	"newsAggr/cmd/types"
	"newsAggr/logger"
)

type JsonParser struct {
	Source string
}

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse() []types.News {
	var news []types.News

	data := extractFileData(sourceToFile[jp.Source])
	if data == nil {
		logger.ErrorLogger.Fatalf("Error extracting file data: %v\n", sourceToFile[jp.Source])
	}

	var tmp types.Json

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}

	news = append(news, tmp.Articles...)

	return news
}
