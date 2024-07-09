package main

import (
	"gogator/cmd"
	"log"
)

func main() {
	rootCmd := cmd.FetchNewsCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
