package cli

import (
	"errors"
	"fmt"
)

const (
	// ExitCodeOK represents a successful CLI command.
	ExitCodeOK = 0
	// ExitCodeFailure represents a generic command failure; detailed semantics stay in JSON error.code.
	ExitCodeFailure = 1
)

// ExitError carries the process exit code a CLI command wants main to use.
type ExitError struct {
	Code int
}

// Error implements error without exposing domain-specific failure details.
func (e ExitError) Error() string {
	return fmt.Sprintf("exit code %d", e.Code)
}

// ExitCode returns the process exit code represented by err.
func ExitCode(err error) int {
	if err == nil {
		return ExitCodeOK
	}

	var exitError ExitError
	if errors.As(err, &exitError) {
		return exitError.Code
	}

	return ExitCodeFailure
}

// IsExitError reports whether err only represents an intended process status.
func IsExitError(err error) bool {
	var exitError ExitError
	return errors.As(err, &exitError)
}
