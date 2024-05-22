package cli

import (
	"context"
	"embed"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

//go:embed tools/structurizr-cli/lib/*.jar tools/structurizr-cli/structurizr.sh tools/structurizr-cli/structurizr.bat
var structurizrCLI embed.FS

const basePath = "/api"

// Config is the primary means to modify the Client
type Config struct {
	BaseURL    *url.URL
	WorkingDir string
	goos       string
}

// Client is the main Structurizr CLI
// Use NewClient to get started
type Client struct {
	config  *Config
	cmdExec CmdExec
}

// NewClient creates a new Structurizr CLI with sensible but overridable defaults
func NewClient(config *Config, executor CmdExec) *Client {
	// Ensuring the operating system of the current platform is set
	config.goos = runtime.GOOS

	return &Client{config, executor}
}

// PushWorkspace push a new version of a workspace from an existing file
func (c *Client) PushWorkspace(
	ctx context.Context,
	id int64,
	key string,
	secret string,
	passphrase string,
	source string,
) error {
	return c.execute(
		ctx,
		"push",
		"-id", strconv.FormatInt(id, 10),
		"-key", key,
		"-secret", secret,
		"-passphrase", passphrase,
		"-workspace", source,
		"-url", c.config.BaseURL.JoinPath(basePath).String(),
		"-merge", "false",
		"-archive", "true",
	)
}

// WorkingDir extracts the embedded Structurizr CLI files to a working directory for easier utilization.
func WorkingDir(ctx context.Context) (string, error) {
	// Get the path to the directory where the executable is running
	dir, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable directory: %v", err)
	}
	dir = filepath.Join(filepath.Dir(dir), "structurizr-cli")

	tflog.Debug(ctx, fmt.Sprintf("Structurizr CLI working directory: %s", dir))

	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}

	err = fs.WalkDir(structurizrCLI, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore the root directory
		if path == "." {
			return nil
		}

		destPath := filepath.Join(
			dir,
			strings.TrimPrefix(strings.TrimPrefix(path, "tools"), "/structurizr-cli"),
		)
		if d.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		data, err := structurizrCLI.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, os.ModePerm)
	})

	if err != nil {
		return "", fmt.Errorf("failed to extract files: %w", err)
	}

	return dir, nil
}

// execute executes the Structurizr CLI commands with provided options on operating systems that support batch or shell scripts.
func (c *Client) execute(ctx context.Context, options ...string) error {
	var name string
	if c.config.goos == "windows" {
		name = filepath.Join(c.config.WorkingDir, "structurizr.bat")
	} else {
		name = filepath.Join(c.config.WorkingDir, "structurizr.sh")
	}

	// Run the command and capture the output
	out, err := c.cmdExec.CombinedOutput(ctx, name, options...)
	if err != nil {
		return fmt.Errorf("error running Structurizr CLI: %v\nOutput: %s", err, string(out))
	}

	tflog.Debug(ctx, fmt.Sprintf("Structurizr CLI output: %s\n", string(out)))

	return nil
}
