package logger

import (
	"log"
	"os"
)

var (
	// ErrorLogger is a logger instance for logging error messages.
	// It writes log messages to the standard output with a specific format:
	// - The log message prefix is "[ERROR:] \t"
	// - The log message includes the short file name, date, and time
	ErrorLogger = log.New(os.Stdout, "[ERROR:] \t", log.Lshortfile|log.Ldate|log.Ltime)

	// InfoLogger is a logger instance for logging informational messages.
	// It writes log messages to the standard output with a specific format:
	// - The log message prefix is "[INFO:] \t"
	// - The log message includes the date and time
	InfoLogger = log.New(os.Stdout, "[INFO:] \t", log.Ldate|log.Ltime)
)
