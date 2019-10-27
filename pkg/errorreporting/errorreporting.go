package errorreporting

import (
	"context"
	"log"

	er "cloud.google.com/go/errorreporting"
	"github.com/majordomusio/commons/pkg/env"
)

var errorClient *er.Client

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

func init() {
	ctx := context.Background()

	c, err := er.NewClient(ctx, env.Getenv("PROJECT_ID", ""), er.Config{
		ServiceName: env.Getenv("SERVICE_NAME", "default"),
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		log.Printf("Could not initialize errorreporting: %v", err)
	}
	errorClient = c
}

// Close releases the error reporting client
func Close() {
	if errorClient == nil {
		return
	}
	errorClient.Flush()
	errorClient.Close()
	errorClient = nil
}

// Report reports the error, what else?
func Report(err error) {
	errorClient.Report(er.Entry{Error: err})
}
