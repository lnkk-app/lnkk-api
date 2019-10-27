package slack

import "fmt"

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// newError returns an error that formats as the given text.
func newError(cmd, text string) error {
	return &errorString{fmt.Sprintf("api='%s',err='%s'", cmd, text)}
}
