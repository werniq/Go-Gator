package main

import (
	"newsAggr/cmd"
	"newsAggr/logger"
)

func main() {
	rootCmd := cmd.InitRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
}
