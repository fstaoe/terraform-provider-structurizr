package client

import (
	"context"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/api/model"
)

type WorkspacesClient interface {
	GetWorkspaces(ctx context.Context) (*model.Workspaces, error)
	CreateWorkspace(ctx context.Context) (*model.Workspace, error)
	DeleteWorkspace(ctx context.Context, id int64) (*model.APIResponse, error)
}

type WorkspaceClient interface {
	// PushWorkspace push a new version of a workspace from an existing file
	PushWorkspace(ctx context.Context, id int64, key string, secret string, passphrase string, source string) error
}

// Manager is managing the required clients to interact with Structurizr.
// Use NewManager to get started
type Manager struct {
	api WorkspacesClient
	cli WorkspaceClient
}

// NewManager creates a new Manager with the required clients to interact with Structurizr
func NewManager(api WorkspacesClient, cli WorkspaceClient) *Manager {
	return &Manager{api, cli}
}

// GetWorkspaces lists all workspaces
func (m *Manager) GetWorkspaces(ctx context.Context) (*model.Workspaces, error) {
	return m.api.GetWorkspaces(ctx)
}

// CreateWorkspace creates a new workspace
func (m *Manager) CreateWorkspace(ctx context.Context) (*model.Workspace, error) {
	return m.api.CreateWorkspace(ctx)
}

// DeleteWorkspace deletes a workspace
func (m *Manager) DeleteWorkspace(ctx context.Context, id int64) (*model.APIResponse, error) {
	return m.api.DeleteWorkspace(ctx, id)
}

// PushWorkspace push a new version of a workspace from an existing file
func (m *Manager) PushWorkspace(
	ctx context.Context,
	id int64,
	key string,
	secret string,
	passphrase string,
	source string,
) error {
	return m.cli.PushWorkspace(ctx, id, key, secret, passphrase, source)
}
