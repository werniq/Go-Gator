package parsers

import (
	"encoding/json"
	"gogator/cmd/types"
	"io"
	"net/http"
)

// JsonParser struct is used to parse articles data from Json files.
//
// It implements Parser interface, which has a method Parse.
// Firstly ot opens file with name from Source field, then decodes its content into
// array of articles, and throws an error if something went wrong.
//
// Returns a successfully decoded array of news, and nil error.
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

	res, err := http.Get(sourceToEndpoint[jp.Source])
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &news)
	if err != nil {
		return nil, err
	}

	return news, nil
}
