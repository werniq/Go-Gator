package parsers

import (
	"bytes"
	"encoding/xml"
	"io"
	"newsAggr/cmd/FilteringInstructions"
	"newsAggr/cmd/types"
	"newsAggr/logger"
	"reflect"
)

type XMLParser struct {
}

// Parse function is required for XMLParser struct, in order to implement NewsParser interface, for data formatted in xml
func (xp XMLParser) Parse(params *types.FilteringParams) []types.News {
	var data []byte
	var err error
	var news []types.News
	b := &bytes.Buffer{}

	var filenames []string

	if !reflect.DeepEqual(params.Sources, filenames) {
		sourceToFile := map[string]string{
			"abc":             "abcnews.xml",
			"bbc":             "bbc.xml",
			"washingtontimes": "washington-times.xml",
		}

		filenameSet := make(map[string]struct{})
		for _, val := range params.Sources {
			if filename, ok := sourceToFile[val]; ok {
				if _, exists := filenameSet[filename]; !exists {
					filenameSet[filename] = struct{}{}
					filenames = append(filenames, filename)
				}
			}
		}

		if len(filenames) == 0 {
			return nil
		}
	} else {
		filenames = []string{"abcnews.xml", "bbc.xml", "washington-times.xml"}
	}

	for _, filename := range filenames {
		var dummy []types.RSS

		b = bytes.NewBuffer([]byte{})
		b.Write(extractFileData(filename))

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

	factory := FilteringInstructions.GoGatorInstructionFactory{}

	return ApplyParams(news, params, factory)
}
