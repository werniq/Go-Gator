package parsers

import (
	"encoding/json"
	"errors"
	"io"
	"newsaggr/cmd/types"
	"os"
	"path/filepath"
)

// AddNewSource inserts new source to available sources list and determines the appropriate Parser for it
func AddNewSource(format, source, endpoint string) error {
	sourceToEndpoint[source] = endpoint
	sourceToParser[source] = determineParser(format, source)

	err := updateSourcesFile()
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
	sourceToEndpoint[source] = newEndpoint
	err := updateSourcesFile()
	if err != nil {
		return err
	}

	return nil
}

// UpdateSourceFormat updates format for the given source
func UpdateSourceFormat(source, format string) error {
	sourceToParser[source] = determineParser(format, source)
	err := updateSourcesFile()
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

		err := updateSourcesFile()
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadSourcesFile initializes sourceToParser and sourceToEndpoint with data from sources.json file
func LoadSourcesFile() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(cwdPath, CmdDir, ParsersDir, DataDir, sourcesFile)

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

// InitSourcesFile is used to initialize file with information about sources
func InitSourcesFile() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(cwdPath, CmdDir, ParsersDir, DataDir, sourcesFile)

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

// updateSourcesFile is used to update file with information about sources to prevent losing all information if server
// crashes
func updateSourcesFile() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(cwdPath, CmdDir, ParsersDir, DataDir, sourcesFile)

	file, err := os.Create(f)
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
