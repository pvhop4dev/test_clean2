package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// InfoLogger logs informational messages
	InfoLogger *log.Logger
	// ErrorLogger logs error messages
	ErrorLogger *log.Logger
)

// Init initializes the logger
func Init(logPath string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Initialize loggers
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if InfoLogger != nil {
		InfoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal logs a fatal error and exits the program
func Fatal(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Output(2, fmt.Sprintf("FATAL: "+format, v...))
	}
	os.Exit(1)
}

// GetCallerInfo returns the file and line number of the caller
func GetCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "", 0
	}
	return file, line
}
