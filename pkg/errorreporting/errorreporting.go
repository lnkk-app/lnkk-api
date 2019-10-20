package errorreporting

import (
	"context"
	"log"

	er "cloud.google.com/go/errorreporting"
	"github.com/majordomusio/commons/pkg/env"
)

var errorClient *er.Client

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
	errorClient.Flush()
	errorClient.Close()
	errorClient = nil
}

// Report reports the error, what else?
func Report(err error) {
	errorClient.Report(er.Entry{Error: err})
}

/*

	Debug(ctx context.Context, topic, format string, args ...interface{})
	Critical(ctx context.Context, topic, format string, args ...interface{})
	Info(ctx context.Context, topic, format string, args ...interface{})
	Warn(ctx context.Context, topic, format string, args ...interface{})
	Error(ctx context.Context, topic, format string, args ...interface{})

*/
