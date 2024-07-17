package parsers

import (
	"encoding/json"
	"gogator/cmd/types"
)

type JsonParser struct {
	Source string
}

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse() ([]types.News, error) {
	var news types.Json

	data, err := extractFileData(sourceToFile[jp.Source])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &news)
	if err != nil {
		return nil, err
	}

	return news.Articles, nil
}
