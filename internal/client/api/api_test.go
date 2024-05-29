package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// MockHTTPClient is a mock for the HTTP doer
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewClient_InsecureSkipVerify(t *testing.T) {
	t.Run("Given disabled TLS verification", func(t *testing.T) {
		client := NewClient(&Config{TLSInsecure: true})
		transport := client.doer.(*http.Client).Transport.(*http.Transport)
		assert.NotNil(t, transport)
		assert.NotNil(t, transport.TLSClientConfig)
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
	t.Run("Given enabled TLS verification", func(t *testing.T) {
		client := NewClient(&Config{TLSInsecure: false})
		transport := client.doer.(*http.Client).Transport.(*http.Transport)
		assert.NotNil(t, transport)
		assert.NotNil(t, transport.TLSClientConfig)
		assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
	t.Run("Given default TLS verification", func(t *testing.T) {
		client := NewClient(&Config{})
		transport := client.doer.(*http.Client).Transport.(*http.Transport)
		assert.NotNil(t, transport)
		assert.NotNil(t, transport.TLSClientConfig)
		assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
}

// TestGetWorkspaces tests the GetWorkspaces function
func TestGetWorkspaces(t *testing.T) {
	config := &Config{
		AdminAPIKey: "test-key",
		BaseURL:     &url.URL{Scheme: "http", Host: "localhost:8080"},
		UserAgent:   "test-agent",
	}

	mockClient := new(MockHTTPClient)
	client := &Client{config, mockClient}

	mockResponse := &model.Workspaces{
		Workspaces: []*model.Workspace{
			{
				ID:   1,
				Name: "Test Workspace",
			},
		},
	}

	body, _ := json.Marshal(mockResponse)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	ctx := context.Background()
	workspaces, err := client.GetWorkspaces(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, workspaces)
	assert.Equal(t, 1, len(workspaces.Workspaces))
	assert.Equal(t, int64(1), workspaces.Workspaces[0].ID)
}

// TestCreateWorkspace tests the CreateWorkspace function
func TestCreateWorkspace(t *testing.T) {
	config := &Config{
		AdminAPIKey: "test-key",
		BaseURL:     &url.URL{Scheme: "http", Host: "localhost:8080"},
		UserAgent:   "test-agent",
	}

	mockClient := new(MockHTTPClient)
	client := &Client{config, mockClient}

	mockResponse := &model.Workspace{
		ID:   1,
		Name: "Test Workspace",
	}

	body, _ := json.Marshal(mockResponse)
	resp := &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	ctx := context.Background()
	workspace, err := client.CreateWorkspace(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, workspace)
	assert.Equal(t, int64(1), workspace.ID)
}

// TestDeleteWorkspace tests the DeleteWorkspace function
func TestDeleteWorkspace(t *testing.T) {
	config := &Config{
		AdminAPIKey: "test-key",
		BaseURL:     &url.URL{Scheme: "http", Host: "localhost:8080"},
		UserAgent:   "test-agent",
	}

	mockClient := new(MockHTTPClient)
	client := &Client{config, mockClient}

	mockResponse := &model.APIResponse{
		Message: "Deleted",
	}

	body, _ := json.Marshal(mockResponse)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	ctx := context.Background()
	apiResponse, err := client.DeleteWorkspace(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, apiResponse)
	assert.Equal(t, "Deleted", apiResponse.Message)
}
