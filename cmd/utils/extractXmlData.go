package utils

import (
	"io"
	"newsAggr/logger"
	"os"
)

// ExtractFileData reads data from file $filename and returns its content
func ExtractFileData(filename string) []byte {
	cwd, err := os.Getwd()
	if err != nil {
		logger.ErrorLogger.Fatalf("Error getting current working directory: %v\n", err)
		return nil
	}
	cwd += "\\data\\"

	file, err := os.Open(cwd + filename)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error opening xml file %s: %v\n", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error reading XML data: %v\n", err)
	}

	return data
}
