package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// SetUpLogger sets up logging to both a file and stdout.
// Parameters:
// - logFileName: The base name for the log file (e.g., "app.log").
// Returns:
// - *os.File: A pointer to the log file opened for writing, so that the caller can close it when done.
// - error: An error is returned if any problem occurs during the setup.
// The log file is created in a "logs" directory with a timestamp prefix.
func LoggerSetUpLoggerFile(logFileName string) (*os.File, error) {
	// Create the "logs" directory if it doesn't already exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get the current timestamp to use as part of the log file name
	currentTime := time.Now().Format("20060102_150405")
	logFileName = fmt.Sprintf("%s_%s", currentTime, logFileName) // Create log file with timestamp prefix

	// Build the full log file path inside the "logs" directory
	logFilePath := "logs/" + logFileName

	// Try to open or create the log file for appending
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// If there's an error opening the log file, fallback to stdout only logging
		log.Printf("Error opening log file %s, falling back to stdout only: %v", logFilePath, err)
		logFile = nil // No log file to use
	} else {
		// Successfully opened the log file
		log.Printf("Logging initialized. Log file: %s", logFilePath)
	}

	// Set up multi-writer to write logs to both stdout and the log file, if available
	var multiWriter io.Writer
	if logFile != nil {
		multiWriter = io.MultiWriter(os.Stdout, logFile) // Log to both stdout and file
	} else {
		multiWriter = os.Stdout // Only log to stdout if the log file couldn't be opened
	}

	// Configure the logger to include timestamp, microseconds, and short file location in the log messages
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	return logFile, nil
}

// logFile, err := logger.SetUpLogger("app.log")
// if err != nil {
// 	log.Fatalf("Failed to set up logger: %v", err)
// }
// defer logFile.Close() // Ensure to close the log file when the application exits
