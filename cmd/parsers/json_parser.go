package parsers

import (
	"encoding/json"
	"newsAggr/cmd/types"
)

type JsonParser struct {
	Source string
}

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse() ([]types.News, error) {
	var news []types.News

	data, err := extractFileData(sourceToFile[jp.Source])
	if data == nil {
		return nil, err
	}

	var tmp types.Json

	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}

	news = append(news, tmp.Articles...)

	return news, nil
}
