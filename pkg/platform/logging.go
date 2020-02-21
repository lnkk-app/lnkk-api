package platform

import (
	"context"
	"fmt"
	"os"
)

// FIXME just stubs for now !

// PrintError prints an error message to stderr
func PrintError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
}

// Debug formats its arguments according to the format, analogous to fmt.Printf, and records the text as a log message at Debug level.
func Debug(topic, format string, args ...interface{}) {
	//logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Debug})
}

// Info is like Debug, but on INFO level
func Info(topic, format string, args ...interface{}) {
	//logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Info})
}

// Warn is like Debug, but on WARN level
func Warn(topic, format string, args ...interface{}) {
	//logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Warning})
}

// Critical is like Debug, but on CRITICAL level
func Critical(topic, format string, args ...interface{}) {
	//logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Warning})
}

// Count collects a numeric counter value
func Count(ctx context.Context, topic, label string, value int) {
	// FIXME just a stub
}
