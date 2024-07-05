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
		log.Println("Error loading .env file: ")
	}

	switch os.Getenv("APP_MODE") {
	case "WINDOWS":
		parsers.PathToDataDir = "\\cmd\\parsers\\data\\"
		server.RelativePathToCertsDir = "\\cmd\\server\\certs\\"
	case "DOCKER":
		parsers.PathToDataDir = "/cmd/parsers/data/"
		server.RelativePathToCertsDir = "/cmd/server/certs/"
	}
}

func main() {
	err := server.ConfAndRun()
	if err != nil {
		log.Fatalln(err)
	}
}
