package platform

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/errorreporting"

	"github.com/majordomusio/commons/pkg/env"
)

var errorClient *errorreporting.Client
var dsClient *datastore.Client

func init() {
	ctx := context.Background()

	// initialize error reporting
	ec, err := errorreporting.NewClient(ctx, env.Getenv("PROJECT_ID", ""), errorreporting.Config{
		ServiceName: env.Getenv("SERVICE_NAME", "default"),
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		log.Printf("Could not initialize errorreporting: %v", err)
	}
	errorClient = ec

	// initialize the datastore
	ds, err := datastore.NewClient(ctx, env.Getenv("PROJECT_ID", ""))
	if err != nil {
		log.Printf("Could not initialize datastore: %v", err)
	}
	dsClient = ds
}

// Close the platform related clients
func Close() {
	// error reporting
	if errorClient == nil {
		return
	}
	errorClient.Flush()
	errorClient.Close()
	errorClient = nil

	// datastore
	if dsClient == nil {
		return
	}
	dsClient.Close()
	dsClient = nil
}

// Report reports the error, what else?
func Report(err error) {
	errorClient.Report(errorreporting.Entry{Error: err})
}

// DataStore return a reference to the datastore client
func DataStore() *datastore.Client {
	return dsClient
}

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
