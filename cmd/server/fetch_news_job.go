package server

import (
	"encoding/json"
	"errors"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"os"
	"path/filepath"
)

type FetchNewsJob struct {
	Filters *types.FilteringParams
}

var (
	// ErrCreatingFile is thrown when there was an error while creating sources file
	ErrCreatingFile = "Error while creating a file: "

	// ErrParsingSource is thrown when we have error while parsing sources
	ErrParsingSource = "Error while parsing sources: "
)

// Run is a function that fetches news, parses it, and writes the parsed data
// to a JSON file named with the current date in the format YYYY-MM-DD.
func (j *FetchNewsJob) Run() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	articleFilepath := filepath.Join(cwdPath,
		parsers.StoragePath,
		j.Filters.StartingTimestamp+parsers.JsonExtension)

	file, err := os.Create(articleFilepath)
	if err != nil {
		return errors.New(ErrCreatingFile + err.Error())
	}

	news, err := parsers.ParseBySource(parsers.AllSources)
	if err != nil {
		return errors.New(ErrParsingSource + err.Error())
	}

	// since some publishers may still store news from previous dates
	// program additionally applies date range filters
	news = filters.Apply(news, j.Filters)

	articlesData, err := json.Marshal(news)
	if err != nil {
		return err
	}

	_, err = file.Write(articlesData)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
