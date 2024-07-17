package main

import (
	"gogator/cmd/server"
	"log"
)

func main() {
	err := server.ConfAndRun()
	if err != nil {
		log.Fatalln("Error configuring and running server: ", err)
	}
}
