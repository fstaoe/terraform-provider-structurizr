package cli

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	client := NewClient(&Config{BaseURL: baseURL, WorkingDir: "/tmp", goos: runtime.GOOS}, DefaultCmdExec)
	expectedClient := &Client{
		config:  &Config{BaseURL: baseURL, WorkingDir: "/tmp", goos: runtime.GOOS},
		cmdExec: DefaultCmdExec,
	}

	assert.IsType(t, expectedClient, client)
	assert.Equal(t, expectedClient, client)
}

func TestPushWorkspace(t *testing.T) {
	cmdExecMock := &mockCmdExec{
		output:       []byte("mocked output"),
		err:          nil,
		expectedName: filepath.Join("/tmp", "structurizr.sh"),
		expectedArgs: []string{
			"push",
			"-id", "12345",
			"-key", "key",
			"-secret", "secret",
			"-passphrase", "passphrase",
			"-workspace", "response_workspace.tmpl",
			"-url", "http://localhost/api",
			"-merge", "false",
			"-archive", "true",
		},
	}
	baseURL, _ := url.Parse("http://localhost")
	client := &Client{config: &Config{baseURL, "/tmp", runtime.GOOS}, cmdExec: cmdExecMock}

	err := client.PushWorkspace(context.TODO(), 12345, "key", "secret", "passphrase", "response_workspace.tmpl")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmdExecMock.capturedName != cmdExecMock.expectedName {
		t.Fatalf("expected command name %q, got %q", cmdExecMock.expectedName, cmdExecMock.capturedName)
	}
	if !reflect.DeepEqual(cmdExecMock.capturedArgs, cmdExecMock.expectedArgs) {
		t.Fatalf("expected command args %v, got %v", cmdExecMock.expectedArgs, cmdExecMock.capturedArgs)
	}
}

func TestExecute(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	type fields struct {
		config  *Config
		cmdExec *mockCmdExec
	}
	type args struct {
		ctx     context.Context
		options []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Given a simple command",
			fields{
				config: &Config{baseURL, "/tmp", runtime.GOOS},
				cmdExec: &mockCmdExec{
					output:       []byte("mocked output"),
					err:          nil,
					expectedName: filepath.Join("/tmp", "structurizr.sh"),
					expectedArgs: []string{"echo"},
				},
			},
			args{
				ctx:     context.TODO(),
				options: []string{"echo"},
			},
			false,
		},
		{
			"Given a command executed on windows",
			fields{
				config: &Config{baseURL, "/tmp", "windows"},
				cmdExec: &mockCmdExec{
					output:       []byte("mocked output"),
					err:          nil,
					expectedName: filepath.Join("/tmp", "structurizr.bat"),
					expectedArgs: []string{"echo"},
				},
			},
			args{
				ctx:     context.TODO(),
				options: []string{"echo"},
			},
			false,
		},
		{
			"Given a failure during command execution",
			fields{
				config: &Config{baseURL, "/tmp", runtime.GOOS},
				cmdExec: &mockCmdExec{
					output:       []byte("mocked output"),
					err:          errors.New("oops, command failed"),
					expectedName: filepath.Join("/tmp", "structurizr.sh"),
					expectedArgs: []string{"echo"},
				},
			},
			args{
				ctx:     context.TODO(),
				options: []string{"echo"},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{config: tt.fields.config, cmdExec: tt.fields.cmdExec}

			if err := c.execute(tt.args.ctx, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.cmdExec.capturedName != tt.fields.cmdExec.expectedName {
				t.Fatalf("expected command name %q, got %q", tt.fields.cmdExec.expectedName, tt.fields.cmdExec.capturedName)
			}
			if !reflect.DeepEqual(tt.fields.cmdExec.capturedArgs, tt.fields.cmdExec.expectedArgs) {
				t.Fatalf("expected command args %v, got %v", tt.fields.cmdExec.expectedArgs, tt.fields.cmdExec.capturedArgs)
			}
		})
	}
}
