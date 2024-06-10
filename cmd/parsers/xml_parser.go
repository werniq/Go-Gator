package parsers

import (
	"bytes"
	"encoding/xml"
	"io"
	"newsAggr/cmd/types"
	"newsAggr/logger"
)

type XMLParser struct {
	Source string
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse() []types.News {
	var news []types.News

	var tmp []types.RSS

	b := bytes.NewBuffer([]byte{})
	b.Write(extractFileData(sourceToFile[xp.Source]))

	data, err := io.ReadAll(b)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error reading data from buffer: %v\n", err)
	}

	err = xml.Unmarshal(data, &tmp)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error decoding XML data: %v\n", err)
	}
	news = append(news, tmp[0].Channel.Items...)

	return news
}
