package parsers

import (
	"encoding/json"
	"gogator/cmd/types"
)

type JsonParser struct {
	Source string
}

// getFileData
//
// This type defines a function which open path by given filename, and returns its
// content, and errors of any occur.
//
// I defined a variable openFile with this type, to mock parsing data in JsonParser
type getFileData func(f string) ([]byte, error)

var openFile getFileData = extractFileData

// Parse function is required for JsonParser struct, in order to implement NewsParser interface, for data formatted in json
func (jp JsonParser) Parse() ([]types.News, error) {
	var news []types.News

	data, err := openFile(jp.Source)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &news)
	if err != nil {
		return nil, err
	}

	return news, nil
}
