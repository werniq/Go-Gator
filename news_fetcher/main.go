package main

import (
	"flag"
	"log"
)

const (
	// defaultStoragePath contains the default path to the directory where all data will be stored
	defaultStoragePath = "/tmp/"
)

func main() {
	var storagePath string

	flag.StringVar(&storagePath, "fs", defaultStoragePath,
		"Path to directory where all data will be stored")
	flag.Parse()

	err := RunJob(storagePath)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully fetched and parsed news")
}
