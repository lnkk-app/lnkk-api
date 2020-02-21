package platform

import (
	"context"
	"log"

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
	if errorClient != nil {
		errorClient.Flush()
		errorClient.Close()
		errorClient = nil
	}

	// datastore
	if dsClient != nil {
		dsClient.Close()
		dsClient = nil
	}
}

// Report reports an error, what else?
func Report(err error) {
	errorClient.Report(errorreporting.Entry{Error: err})
}

// DataStore returns a reference to the datastore
func DataStore() *datastore.Client {
	return dsClient
}
