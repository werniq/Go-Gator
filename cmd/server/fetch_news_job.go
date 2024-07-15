package server

import (
	"encoding/json"
	"errors"
	"newsaggr/cmd/filters"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
	"os"
	"path/filepath"
)

type FetchNewsJob struct {
	Filters *types.FilteringParams
}

var (
	ErrCreatingFile  = "Error while creating a file: "
	ErrParsingSource = "Error while parsing sources: "
)

// Run is a function that fetches news, parses it, and writes the parsed data
// to a JSON file named with the current date in the format YYYY-MM-DD.
func (j *FetchNewsJob) Run() error {
	CwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(CwdPath,
		parsers.CmdDir, parsers.ParsersDir, parsers.DataDir,
		j.Filters.StartingTimestamp+".json")

	file, err := os.Create(f)
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

	data, err := json.Marshal(news)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
