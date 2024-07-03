package main

import (
	"github.com/joho/godotenv"
	"log"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/server"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file: ")
	}

	switch os.Getenv("APP_MODE") {
	case "DEVELOPMENT":
		parsers.PathToDataDir = "\\cmd\\parsers\\data\\"
	case "DOCKER":
		parsers.PathToDataDir = "/cmd/parsers/data/"
	}
}

func main() {
	server.ConfAndRun()
}
