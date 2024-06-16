package main

import (
	"log"
	"newsaggr/cmd"
)

func main() {
	rootCmd := cmd.InitRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
