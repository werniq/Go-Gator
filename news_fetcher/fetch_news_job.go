package main

import (
	"encoding/json"
	"errors"
	"gogator/cmd/filters"
	"gogator/cmd/parsers"
	"gogator/cmd/types"
	"log"
	"os"
	"path/filepath"
	"time"
)

// NewsFetchingJob struct is used to fetch and parse articles feeds,
// and then writes the parsed data to a JSON file named with the current date
//
// Using Kubernetes CronJob object, it will run once in a day, to parse
type NewsFetchingJob struct {
	params      *types.FilteringParams
	storagePath string
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

// RunJob initializes and runs NewsFetchingJob, which will parse data from feeds into respective files
func RunJob(storagePath string) error {
	dateTimestamp := time.Now().Format(time.DateOnly)
	job := &NewsFetchingJob{
		params: types.NewFilteringParams("",
			dateTimestamp,
			"",
			""),
		storagePath: storagePath,
	}

	err := job.Execute()
	if err != nil {
		return err
	}

	return nil
}

// Execute is a function that fetches news, parses it, and writes the parsed data
// to a JSON file named with the current date in the format YYYY-MM-DD.
func (j *NewsFetchingJob) Execute() error {
	articleFilepath := filepath.Join(
		j.storagePath,
		j.params.StartingTimestamp+".json")

	articlesFile, err := os.Create(articleFilepath)
	if err != nil {
		return errors.New(errCreatingFile + err.Error())
	}

	defer func(articlesFile *os.File) {
		err = articlesFile.Close()
		if err != nil {
			log.Fatalln(errClosingFile + err.Error())
		}
	}(articlesFile)

	news, err := parsers.ParseBySource(parsers.AllSources)
	if err != nil {
		return errors.New(errParsingSources + err.Error())
	}

	news = filters.Apply(news, j.params)

	articlesData, err := json.Marshal(news)
	if err != nil {
		return errors.New(errMarshalData + err.Error())
	}

	_, err = articlesFile.Write(articlesData)
	if err != nil {
		return errors.New(errWritingData + err.Error())
	}

	return nil
}
