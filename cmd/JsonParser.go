package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"newsAggr/cmd/types"
	"newsAggr/cmd/utils"
	"newsAggr/logger"
)

type JSONParser struct {
}

// Parse function is required for JSONParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JSONParser) Parse(params ParsingParams) []types.News {
	var data []byte
	var err error
	var news []types.News
	b := &bytes.Buffer{}

	filenames := []string{"nbc-news.json"}

	if !utils.ArrayIncludes(params.Sources, filenames[0]) {
		return nil
	}

	for _, filename := range filenames {
		var dummy types.Json

		b = bytes.NewBuffer([]byte{})
		b.Write(utils.ExtractFileData(filename))

		data, err = io.ReadAll(b)
		if err != nil {
			logger.ErrorLogger.Fatalf("Error reading data from buffer: %v\n", err)
		}

		err = json.Unmarshal(data, &dummy)
		if err != nil {
			logger.ErrorLogger.Fatalf("Error decoding JSON data: %v\n", err)
		}

		news = append(news, types.JsonNewsToNews(dummy.Articles)...)
	}

	factory := GoGatorInstructionFactory{}

	return ApplyParams(news, params, factory)
}
