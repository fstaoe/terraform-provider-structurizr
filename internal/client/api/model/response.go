package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// APIErrBadRequest represents an HTTP 400 error
	APIErrBadRequest = errors.New("bad request")
	// APIErrSystemUnavailable represents an hTTP 5xx error
	APIErrSystemUnavailable = errors.New("system unavailable")
	// APIErrUnauthorized represents an HTTP 401 error
	APIErrUnauthorized = errors.New("unauthorized")
)

// APIResponse represents a response from structurizr
type APIResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Revision int64  `json:"revision"`
}

// APIErrorResponse represents a body of the shape: {"success":false,"message":"error message"}
type APIErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Err     error
}

// Error implements the error interface
func (e *APIErrorResponse) Error() string {
	message := new(strings.Builder)
	if e.Message != "" {
		message.WriteString("\terror details: \n")
		message.WriteString(fmt.Sprintf("\t\tsummary: %s\n", e.Message))
	} else {
		message.WriteString("\n\tplease see the log for error details\n")
	}

	if e.Err != nil {
		return fmt.Sprintf("%s \n\n%s", e.Err, message.String())
	}

	return message.String()
}
