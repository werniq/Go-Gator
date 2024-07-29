package fetchnews

import (
	"encoding/json"
	"errors"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/server/handlers"
	"gogator/cmd/types"
	"os"
	"path/filepath"
	"time"
)

// FetchingJob struct is used to fetch and parse articles feeds,
// and then writes the parsed data to a JSON file named with the current date
//
// Using Kubernetes CronJob object, it will run once in a day, to parse
type FetchingJob struct {
	filters *types.FilteringParams
}

const (
	// errCreatingFile is thrown when there was an error while creating sources file
	errCreatingFile = "Error while creating a file: "

	// errParsingSources is thrown when we have error while parsing sources
	errParsingSources = "Error while parsing sources: "

	// errMarshalData is used for better error logging when we have error during json.Marshal call
	errMarshalData = "Error while performing JSON encoding articles: "

	// errWritingData is thrown when we have error during writing data to the file
	errWritingData = "Error while writing data to file: "

	// errClosingFile is thrown when we have error while closing file
	errClosingFile = "Error closing file: "
)

// RunJob initializes and runs FetchingJob, which will parse data from feeds into respective files
func RunJob() error {
	dateTimestamp := time.Now().Format(time.DateOnly)
	j := FetchingJob{
		filters: types.NewFilteringParams("",
			dateTimestamp,
			"",
			""),
	}

	err := j.Run()
	if err != nil {
		return err
	}

	handlers.LastFetchedFileDate = dateTimestamp

	return nil
}

// Run is a function that fetches news, parses it, and writes the parsed data
// to a JSON file named with the current date in the format YYYY-MM-DD.
func (j *FetchingJob) Run() error {
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	articleFilepath := filepath.Join(cwdPath,
		parsers.StoragePath,
		j.filters.StartingTimestamp+".json")

	file, err := os.Create(articleFilepath)
	if err != nil {
		return errors.New(errCreatingFile + err.Error())
	}

	news, err := parsers.ParseBySource(parsers.AllSources)
	if err != nil {
		return errors.New(errParsingSources + err.Error())
	}

	news = filters.Apply(news, j.filters)

	articlesData, err := json.Marshal(news)
	if err != nil {
		return errors.New(errMarshalData + err.Error())
	}

	_, err = file.Write(articlesData)
	if err != nil {
		return errors.New(errWritingData + err.Error())
	}

	err = file.Close()
	if err != nil {
		return errors.New(errClosingFile + err.Error())
	}

	return nil
}
