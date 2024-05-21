package main

<<<<<<< HEAD
func main() {
  
=======
import (
	"newsAggr/cmd"
	"newsAggr/logger"
)

func init() {

}

func main() {
	rootCmd := cmd.InitRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
>>>>>>> 579d012 (Go-Gator: version: 1.0 | Filtering news by various options)
}
