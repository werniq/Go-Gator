package main

import (
	"log"
	"newsaggr/cmd/server"
)

func main() {
	err := server.ConfAndRun()
	if err != nil {
		log.Fatalln(err)
	}
}
