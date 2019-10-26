package store

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/majordomusio/commons/pkg/env"
)

var dsClient *datastore.Client

func init() {
	ctx := context.Background()
	c, err := datastore.NewClient(ctx, env.Getenv("PROJECT_ID", ""))
	if err != nil {
		errorreporting.Report(err)
	}
	dsClient = c
}

// Client return a reference to the datastore client
func Client() *datastore.Client {
	return dsClient
}

// Close does the clean-up
func Close() {
	if dsClient == nil {
		return
	}
	dsClient.Close()
	dsClient = nil
}
