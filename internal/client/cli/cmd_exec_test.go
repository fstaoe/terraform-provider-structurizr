package cli

import (
	"context"
	"reflect"
	"testing"
)

// mockCmdExec is a mock implementation of CmdExec for testing
type mockCmdExec struct {
	expectedName string
	expectedArgs []string
	output       []byte
	err          error
	capturedName string
	capturedArgs []string
}

// CombinedOutput is capturing and storing the input so later it can be asserted
func (m *mockCmdExec) CombinedOutput(_ context.Context, name string, arg ...string) ([]byte, error) {
	m.capturedName = name
	m.capturedArgs = arg
	return m.output, m.err
}

func TestCmdExec_CombinedOutput_Simple(t *testing.T) {
	ctx := context.Background()
	executor := DefaultCmdExec

	// Test a simple command
	output, err := executor.CombinedOutput(ctx, "echo", "hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOutput := "hello\n"
	if string(output) != expectedOutput {
		t.Fatalf("expected output %q, got %q", expectedOutput, string(output))
	}

	// Test a command that returns an error
	_, err = executor.CombinedOutput(ctx, "false")
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
}

func TestCmdExec_CombinedOutput(t *testing.T) {
	mockExecutor := &mockCmdExec{
		output:       []byte("mocked output"),
		err:          nil,
		expectedName: "test-command",
		expectedArgs: []string{"arg1", "arg2"},
	}

	ctx := context.Background()
	output, err := mockExecutor.CombinedOutput(ctx, mockExecutor.expectedName, []string{"arg1", "arg2"}...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(output) != string(mockExecutor.output) {
		t.Fatalf("expected output %q, got %q", string(mockExecutor.output), string(output))
	}

	if mockExecutor.capturedName != mockExecutor.expectedName {
		t.Fatalf("expected command name %q, got %q", mockExecutor.expectedName, mockExecutor.capturedName)
	}

	if !reflect.DeepEqual(mockExecutor.capturedArgs, mockExecutor.expectedArgs) {
		t.Fatalf("expected command args %v, got %v", mockExecutor.expectedArgs, mockExecutor.capturedArgs)
	}
}
