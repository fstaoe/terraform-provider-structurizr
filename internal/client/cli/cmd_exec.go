package cli

import (
	"context"
	"os/exec"
)

// CmdExec is an interface for executing commands
type CmdExec interface {
	CombinedOutput(ctx context.Context, name string, arg ...string) ([]byte, error)
}

// cmdExecutor is an implementation of cmdExecutor that uses exec.Command
type cmdExecutor struct{}

// DefaultCmdExec is the default cmdExecutor.
var DefaultCmdExec = &cmdExecutor{}

// CombinedOutput runs a command and returns its combined standard output and standard error
func (e *cmdExecutor) CombinedOutput(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	return cmd.CombinedOutput()
}
