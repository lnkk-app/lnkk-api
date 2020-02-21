package errors

//
// This is a copy of majordomusio/commons/pkg/error, with improvements.
// Will be eventually backported.
//

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type (
	basicError struct {
		err  error
		msg  string
		oper string
		pkg  string
		fn   string
	}
)

func (e *basicError) Error() string {
	return e.msg
}

func (e *basicError) Unwrap() error {
	return e.err
}

func (e *basicError) FullError() string {
	return fmt.Sprintf("%s. Operation=%s, package='%s',f='%s'", e.msg, e.oper, e.pkg, e.fn)
}

// New returns an error that formats as the given text
func New(text string) error {
	p, f := packageAndFunc()
	return &basicError{err: errors.New(text), msg: text, oper: "undefined", pkg: p, fn: f}
}

// NewError wraps an error with additional metadata
func NewError(e error) error {
	p, f := packageAndFunc()
	return &basicError{err: e, msg: e.Error(), oper: "undefined", pkg: p, fn: f}
}

// NewOperationError wraps an error with additional metadata
func NewOperationError(operation string, e error) error {
	p, f := packageAndFunc()
	return &basicError{err: e, msg: e.Error(), oper: operation, pkg: p, fn: f}
}

// see https://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func packageAndFunc() (string, string) {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkg := ""
	fn := parts[pl-1]
	if parts[pl-2][0] == '(' {
		fn = parts[pl-2] + "." + fn
		pkg = strings.Join(parts[0:pl-2], ".")
	} else {
		pkg = strings.Join(parts[0:pl-1], ".")
	}
	return pkg, fn
}
