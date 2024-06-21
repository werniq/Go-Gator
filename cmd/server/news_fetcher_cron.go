package server

import (
	"encoding/json"
	"fmt"
	"log"
	"newsaggr/cmd/parsers"
	"os"
	"time"
)

func FetchNewsJob() error {
	currentTimestamp := time.Now().Format(time.DateOnly)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s\\cmd\\parsers\\data\\%s.json", wd, currentTimestamp)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if file == nil {
		log.Fatalln("File is nil")
	}

	news, err := parsers.ParseBySource("")
	if err != nil {
		return err
	}

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
