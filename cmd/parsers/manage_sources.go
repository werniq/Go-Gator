package parsers

import (
	"encoding/json"
	"errors"
	"gogator/cmd/types"
	"io"
	"os"
	"path/filepath"
)

const (
	ErrNoSource = "no source was detected. please, create source first"
)

// AddNewSource inserts new source to available sources list and determines the appropriate Parser for it
func AddNewSource(format, source, endpoint string) error {
	sourceToEndpoint[source] = endpoint
	sourceToParser[source] = determineParser(format, source)

	err := UpdateSourcesFile()
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
func UpdateSourceEndpoint(source, newEndpoint string) error {
	if _, exists := sourceToParser[source]; exists {
		return errors.New(ErrNoSource)
	}

	sourceToEndpoint[source] = newEndpoint
	err := UpdateSourcesFile()
	if err != nil {
		return err
	}

	return nil
}

// UpdateSourceFormat updates format for the given source
func UpdateSourceFormat(source, format string) error {
	if _, exists := sourceToParser[source]; exists {
		return errors.New(ErrNoSource)
	}

	sourceToParser[source] = determineParser(format, source)
	err := UpdateSourcesFile()
	if err != nil {
		return err
	}

	return nil
}

// DeleteSource removes source from the map
func DeleteSource(source string) error {
	if _, exists := sourceToEndpoint[source]; exists {
		delete(sourceToEndpoint, source)
		delete(sourceToParser, source)

		err := UpdateSourcesFile()
		if err != nil {
			return err
		}
	} else {
		return errors.New(ErrNoSource)
	}
	return nil
}

// LoadSourcesFile initializes sourceToParser and sourceToEndpoint with data from sources.json file
func LoadSourcesFile() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(cwdPath, StoragePath, sourcesFile)

	file, err := os.Open(f)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	var sources []types.Source
	err = json.Unmarshal(data, &sources)
	if err != nil {
		return err
	}

	for _, s := range sources {
		sourceToParser[s.Name] = determineParser(s.Format, s.Name)
		sourceToEndpoint[s.Name] = s.Endpoint
	}

	return nil
}

// UpdateSourcesFile is used to initialize file with information about sources
func UpdateSourcesFile() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(cwdPath, StoragePath, sourcesFile)

	file, err := os.Create(f)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrExist):
			file, err = os.Open(f)
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
