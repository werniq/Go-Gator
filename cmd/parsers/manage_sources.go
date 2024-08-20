package parsers

import (
	"encoding/json"
	"errors"
	"gogator/cmd/types"
	"io"
	"os"
	"path/filepath"
)

var (
	// g is Parsing factory.
	// These are custom types which will be used for parsers initialization for
	// different data formats
	g ParsingFactory

	// StoragePath is the path to folder with all data from application
	StoragePath string = "/tmp/"

	// sourceToEndpoint maps source names (as strings) to their corresponding filenames
	sourceToEndpoint = map[string]string{
		WashingtonTimes: "https://www.washingtontimes.com/rss/headlines/news/world",
		ABC:             "https://abcnews.go.com/abcnews/internationalheadlines",
		BBC:             "https://feeds.bbci.co.uk/news/rss.xml",
		UsaToday:        "https://usatoday.com",
	}

	// sourceToParser maps source names to their corresponding parser instances
	// The parser are created using ParsingFactory
	sourceToParser = map[string]Parser{
		UsaToday:        g.HtmlParser(UsaToday),
		ABC:             g.XmlParser(ABC),
		BBC:             g.XmlParser(BBC),
		WashingtonTimes: g.XmlParser(WashingtonTimes),
	}
)

// AddNewSource inserts new source to available sources list and determines the appropriate Parser for it
//
// Throws an error, if the source was already registered previously.
func AddNewSource(format, source, endpoint string) error {
	sourceToEndpoint[source] = endpoint
	sourceToParser[source] = determineParser(format, source)

	err := UpdateSourceFile()
	if err != nil {
		return err
	}

	return nil
}

// GetAllSources returns all available sources
func GetAllSources() map[string]string {
	return sourceToEndpoint
}

// GetSourceDetailed returns detailed information about source
func GetSourceDetailed(source string) types.Source {
	return types.Source{
		Name:     source,
		Format:   determineFormat(sourceToParser[source], source),
		Endpoint: sourceToEndpoint[source],
	}
}

// UpdateSourceEndpoint updates endpoint for the given source
//
// Throws an error, if provided source not exists
func UpdateSourceEndpoint(source, newEndpoint string) error {
	sourceToEndpoint[source] = newEndpoint
	err := UpdateSourceFile()
	if err != nil {
		return err
	}

	return nil
}

// UpdateSourceFormat updates format for the given source
//
// Throws an error, if provided source not exists
func UpdateSourceFormat(source, format string) error {
	sourceToParser[source] = determineParser(format, source)
	err := UpdateSourceFile()
	if err != nil {
		return err
	}

	return nil
}

// DeleteSource removes source from the map
func DeleteSource(source string) error {
	delete(sourceToEndpoint, source)
	delete(sourceToParser, source)

	err := UpdateSourceFile()
	if err != nil {
		return err
	}

	return nil
}

// LoadSourcesFile initializes sourceToParser and sourceToEndpoint with data stored in
// sources.json file.
func LoadSourcesFile() error {
	sourcesFilepath := filepath.Join(StoragePath, sourcesFile)

	file, err := os.Open(sourcesFilepath)
	if err != nil {
		return err
	}

	sourcesFileData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if sourcesFileData == nil {
		return nil
	}

	var sources []types.Source
	err = json.Unmarshal(sourcesFileData, &sources)
	if err != nil {
		return err
	}

	for _, s := range sources {
		sourceToParser[s.Name] = determineParser(s.Format, s.Name)
		sourceToEndpoint[s.Name] = s.Endpoint
	}

	return nil
}

// UpdateSourceFile initializes or updates a file with all information about sources.
// It creates the file if it doesn't exist, and updates its content if it does.
//
// Returns an error if the current working directory cannot be retrieved,
// the file cannot be created or opened,
// or if the file content cannot be written or closed properly.
func UpdateSourceFile() error {
	sourcesFilePath := filepath.Join(StoragePath, sourcesFile)

	file, err := os.Create(sourcesFilePath)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrExist):
			file, err = os.Open(sourcesFilePath)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		}

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

	sourcesFileData, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	_, err = file.Write(sourcesFileData)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

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
