// Package runerr standardizes exit codes for command errors.
package runerr

import "errors"

// ExitError wraps an error with a suggested exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string { return e.Err.Error() }

func (e ExitError) Unwrap() error { return e.Err }

// WithCode wraps err with the provided exit code.
func WithCode(err error, code int) error {
	return ExitError{Code: code, Err: err}
}

// Code extracts a code from err if present, otherwise returns defaultCode.
func Code(err error, defaultCode int) int {
	var exitErr ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	return defaultCode
}
