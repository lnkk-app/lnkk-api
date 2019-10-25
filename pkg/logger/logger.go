package logger

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/appengine/log"
)

// PrintError prints an error message to stderr
func PrintError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
}

// Debug formats its arguments according to the format, analogous to fmt.Printf, and records the text as a log message at Debug level.
func Debug(ctx context.Context, topic, format string, args ...interface{}) {
	if len(args) == 0 {
		log.Debugf(ctx, "%s: %s", topic, format)
	} else {
		log.Debugf(ctx, topic+": "+format, args...)
	}
}

// Info is like Debug, but on INFO level
func Info(ctx context.Context, topic, format string, args ...interface{}) {
	if len(args) == 0 {
		log.Infof(ctx, "%s: %s", topic, format)
	} else {
		log.Infof(ctx, topic+": "+format, args...)
	}
}

// Warn is like Debug, but on WARN level
func Warn(ctx context.Context, topic, format string, args ...interface{}) {
	if len(args) == 0 {
		log.Warningf(ctx, "%s: %s", topic, format)
	} else {
		log.Warningf(ctx, topic+": "+format, args...)
	}
}

// Critical is like Debug, but on CRITICAL level
func Critical(ctx context.Context, topic, format string, args ...interface{}) {
	if len(args) == 0 {
		log.Criticalf(ctx, "%s: %s", topic, format)
	} else {
		log.Criticalf(ctx, topic+": "+format, args...)
	}
}

// Error is like Debug, but on ERROR level
// FIXME: how does this fit with errorreporting.Report ?
func Error(ctx context.Context, topic, format string, args ...interface{}) {
	if len(args) == 0 {
		log.Errorf(ctx, "%s: %s", topic, format)
	} else {
		log.Errorf(ctx, topic+": "+format, args...)
	}
}
