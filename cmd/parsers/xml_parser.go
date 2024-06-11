package parsers

import (
	"bytes"
	"encoding/xml"
	"io"
	"newsAggr/cmd/types"
)

type XMLParser struct {
	Source string
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse() ([]types.News, error) {
	var news []types.News

	var tmp []types.RSS

	b := bytes.NewBuffer([]byte{})
	data, err := extractFileData(sourceToFile[xp.Source])
	if err != nil {
		return nil, err
	}
	b.Write(data)

	data, err = io.ReadAll(b)
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
