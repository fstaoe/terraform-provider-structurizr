package model

import (
	"fmt"
	"testing"
)

func TestAPIErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		response APIErrorResponse
		expected string
	}{
		{
			name: "Given message and error",
			response: APIErrorResponse{
				Message: "An error occurred",
				Err:     fmt.Errorf("database connection failed"),
			},
			expected: "database connection failed \n\n\terror details: \n\t\tsummary: An error occurred\n",
		},
		{
			name: "Given message and no error",
			response: APIErrorResponse{
				Message: "An error occurred",
				Err:     nil,
			},
			expected: "\terror details: \n\t\tsummary: An error occurred\n",
		},
		{
			name: "Given no message and error",
			response: APIErrorResponse{
				Message: "",
				Err:     fmt.Errorf("database connection failed"),
			},
			expected: "database connection failed \n\n\n\tplease see the log for error details\n",
		},
		{
			name: "Given no message and no error",
			response: APIErrorResponse{
				Message: "",
				Err:     nil,
			},
			expected: "\n\tplease see the log for error details\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.response.Error()
			if actual != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, actual)
			}
		})
	}
}
