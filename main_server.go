package main

import (
	"github.com/joho/godotenv"
	"log"
	"newsaggr/cmd/server"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file: ")
	}
}

func main() {
	server.ConfAndRun()
}
