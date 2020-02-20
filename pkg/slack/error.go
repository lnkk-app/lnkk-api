package slack

import (
	"errors"
	"fmt"
)

// slackError is a trivial implementation of error.
type slackError struct {
	api string
	err error
	msg string
}

func (e *slackError) Error() string {
	return fmt.Sprintf("api='%s',err='%s'", e.api, e.msg)
}

func (e *slackError) Unwrap() error {
	return e.err
}

// NewError returns an error that formats as the given text.
func NewError(api string, e error) error {
	return &slackError{api: api, err: e, msg: e.Error()}
}

// NewSimpleError returns an error that formats as the given text.
func NewSimpleError(api, text string) error {
	return &slackError{api: api, err: errors.New(text), msg: text}
}
