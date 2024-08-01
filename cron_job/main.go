package main

import (
	"flag"
	"gogator/cron_job/news_fetcher"
	"log"
)

const (
	// defaultStoragePath contains the default path to the directory where all data will be stored
	defaultStoragePath = "/tmp/data"
)

func main() {
	var storagePath string

	flag.StringVar(&storagePath, "fs", defaultStoragePath,
		"Path to directory where all data will be stored")
	flag.Parse()

	err := news_fetcher.RunJob(storagePath)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully fetched and parsed news")
}
