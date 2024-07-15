package main

import (
	"github.com/joho/godotenv"
	"log"
	"newsaggr/cmd/server"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}
}

func main() {
	err := server.ConfAndRun()
	if err != nil {
		log.Fatalln("Error configuring and running server: ", err)
	}
}
