package parsers

import (
	"encoding/xml"
	"newsAggr/cmd/types"
)

type XMLParser struct {
	Source string
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse() ([]types.News, error) {
	var news []types.News

	var tmp []types.RSS

	data, err := extractFileData(sourceToFile[xp.Source])
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}
	news = append(news, tmp[0].Channel.Items...)

	return news, nil
}
