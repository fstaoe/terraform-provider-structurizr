package client

import (
	"errors"
	"fmt"
	"strings"
)

// apiErrorResponse represents a body of the shape: {"success":false,"message":"error message"}
type apiErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	err     error
}

// Error implements the error interface
func (e *apiErrorResponse) Error() string {
	message := new(strings.Builder)
	if e.Message != "" {
		message.WriteString("\terror details: \n")
		message.WriteString(fmt.Sprintf("\t\tsummary: %s\n", e.Message))
	} else {
		message.WriteString("\n\tplease see the log for error details\n")
	}

	if e.err != nil {
		return fmt.Sprintf("%s \n\n%s", e.err, message.String())
	}

	return message.String()
}

var (
	// ErrBadRequest represents an HTTP 400 error
	ErrBadRequest = errors.New("bad request")
	// ErrSystemUnavailable represents an hTTP 5xx error
	ErrSystemUnavailable = errors.New("system unavailable")
	// ErrUnauthorized represents an HTTP 401 error
	ErrUnauthorized = errors.New("unauthorized")
)
