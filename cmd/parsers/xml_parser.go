package parsers

import (
	"encoding/xml"
	"newsaggr/cmd/types"
)

type XMLParser struct {
	Source string
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse() ([]types.News, error) {
	var news []types.RSS

	data, err := extractFileData(sourceToFile[xp.Source])
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &news)
	if err != nil {
		return nil, err
	}

	return news[0].Channel.Items, nil
}
