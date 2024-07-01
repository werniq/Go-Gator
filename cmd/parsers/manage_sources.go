package parsers

import (
	"encoding/json"
	"fmt"
	"newsaggr/cmd/types"
	"os"
)

// determineParser is used to determine which parser we will need for that source
func determineParser(format, source string) Parser {
	switch format {
	case "json":
		return JsonParser{Source: source}
	case "xml":
		return XMLParser{Source: source}
	case "html":
		return HtmlParser{Source: source}
	}
	return nil
}

// determineParser is used to determine which format we will use based on the Parser
func determineFormat(p Parser, source string) string {
	switch p {
	case JsonParser{Source: source}:
		return "json"
	case XMLParser{Source: source}:
		return "xml"
	case HtmlParser{Source: source}:
		return "html"
	}
	return ""
}

// AddNewSource inserts new source to available sources list and determines the appropriate Parser for it
func AddNewSource(format, source, endpoint string) {
	sourceToEndpoint[source] = endpoint
	sourceToParser[source] = determineParser(format, source)
}

// GetAllSources returns all available sources
func GetAllSources() map[string]string {
	return sourceToEndpoint
}

// UpdateSourceEndpoint updates endpoint for the given source
func UpdateSourceEndpoint(source, newEndpoint string) {
	sourceToEndpoint[source] = newEndpoint
}

// UpdateSourceFormat updates format for the given source
func UpdateSourceFormat(source, format string) {
	sourceToParser[source] = determineParser(format, source)
}

// DeleteSource removes source from the map
func DeleteSource(source string) {
	if _, exists := sourceToEndpoint[source]; exists {
		sourceToEndpoint[source] = ""
		sourceToParser[source] = nil
	}
}

// updateSourcesFile is used to update file with information about sources to prevent losing all information if server
// crashes
func updateSourcesFile() error {
	const sourcesFile = "sources.json"

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s%s%s.json", wd, PathToSourcesFile, sourcesFile)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	var sources []types.Source
	for key, val := range sourceToEndpoint {
		sources = append(sources, types.Source{
			Name:     key,
			Format:   determineFormat(sourceToParser[key], key),
			Endpoint: val,
		})
	}

	out, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	_, err = file.Write(out)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
