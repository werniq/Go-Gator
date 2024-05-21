package cmd

import (
	"bytes"
	"encoding/xml"
	"io"
	"newsAggr/cmd/types"
	"newsAggr/cmd/utils"
	"newsAggr/logger"
)

type XMLParser struct {
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse(params ParsingParams) []types.News {
	var data []byte
	var err error
	var news []types.News
	b := &bytes.Buffer{}

	filenames := []string{"abcnews.xml", "bbc.xml", "washington-times.xml"}

	for _, filename := range filenames {
		var dummy []types.RSS

		b = bytes.NewBuffer([]byte{})
		b.Write(utils.ExtractFileData(filename))

		data, err = io.ReadAll(b)
		if err != nil {
			logger.ErrorLogger.Fatalf("Error reading data from buffer: %v\n", err)
		}

		err = xml.Unmarshal(data, &dummy)
		if err != nil {
			logger.ErrorLogger.Fatalf("Error decoding XML data: %v\n", err)
		}
		news = append(news, dummy[0].Channel.Items...)
	}

	factory := GoGatorInstructionFactory{}

	return ApplyParams(news, params, factory)
}
