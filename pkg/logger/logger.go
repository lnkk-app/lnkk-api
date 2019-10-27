package logger

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/majordomusio/commons/pkg/env"
)

var logClient *logging.Client
var logger *logging.Logger

type (
	// LogEntry is the custom logging entity
	LogEntry struct {
		Topic   string
		Message string
	}
)

func init() {
	ctx := context.Background()
	l, err := logging.NewClient(ctx, env.Getenv("PROJECT_ID", ""))
	if err != nil {
		errorreporting.Report(err)
	}
	logClient = l
	logger = logClient.Logger(env.Getenv("PROJECT_ID", ""))
}

// Close closes the connection to the Stackdriver server
func Close() {
	if logClient == nil {
		return
	}
	logClient.Close()
	logClient = nil
	logger = nil
}

func createLogEntry(topic, format string, args ...interface{}) *LogEntry {
	var e LogEntry
	if len(args) == 0 {
		e = LogEntry{
			Topic:   topic,
			Message: format,
		}
	} else {
		e = LogEntry{
			Topic:   topic,
			Message: fmt.Sprintf(format, args...),
		}
	}
	return &e
}

// PrintError prints an error message to stderr
func PrintError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
}

// Debug formats its arguments according to the format, analogous to fmt.Printf, and records the text as a log message at Debug level.
func Debug(topic, format string, args ...interface{}) {
	logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Debug})
}

// Info is like Debug, but on INFO level
func Info(topic, format string, args ...interface{}) {
	logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Info})
}

// Warn is like Debug, but on WARN level
func Warn(topic, format string, args ...interface{}) {
	logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Warning})
}

// Critical is like Debug, but on CRITICAL level
func Critical(topic, format string, args ...interface{}) {
	logger.Log(logging.Entry{Payload: createLogEntry(topic, format, args...), Severity: logging.Warning})
}
