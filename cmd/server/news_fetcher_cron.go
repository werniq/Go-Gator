package server

import (
	"encoding/json"
	"fmt"
	"log"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/types"
	"os"
	"time"
)

const (
	PathToDataDir = "cmd\\parsers\\data"
)

// FetchNewsJob is a function that fetches news, parses it, and writes the parsed data
// to a JSON file named with the current date in the format YYYY-MM-DD.
func FetchNewsJob() error {
	dateTimestamp := time.Now().Format(time.DateOnly)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s\\%s\\%s.json", wd, PathToDataDir, dateTimestamp)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if file == nil {
		log.Fatalln("File is nil")
	}

	news, err := parsers.ParseBySource(parsers.AllSources)
	if err != nil {
		return err
	}

	filters := types.NewFilteringParams("", dateTimestamp, "")
	news = parsers.ApplyFilters(news, filters)

	// marshalling data into json format
	data, err := json.Marshal(news)
	if err != nil {
		return err
	}

	// writing json data
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
