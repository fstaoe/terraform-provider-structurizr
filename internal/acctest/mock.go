package acctest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockEndpoint represents a basic request and response that can be used for creating simple httptest server routes.
type MockEndpoint struct {
	Request  *MockRequest
	Response *MockResponse
	Calls    int
	called   int
}

// MockRequest represents a basic HTTP request
type MockRequest struct {
	Method string
	Uri    string
	Body   *string
}

// MockResponse represents a basic HTTP response.
type MockResponse struct {
	StatusCode  int
	Body        string
	ContentType string
}

// NewMockServer establishes a httptest server to simulate behaviour of Structurizr Server
func NewMockServer(t *testing.T, name string, endpoints []*MockEndpoint) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			http.Error(w, fmt.Sprintf("Error reading from HTTP Request Body: %s", err), http.StatusInternalServerError)
			return
		}
		requestBody := buf.String()

		t.Logf("[DEBUG] Received %s API with: %q %q: %s", name, r.Method, r.RequestURI, requestBody)

		for _, e := range endpoints {
			if r.Method == e.Request.Method && r.RequestURI == e.Request.Uri && (e.Request.Body == nil || requestBody == *e.Request.Body) {
				t.Logf("[DEBUG] Respond %s API with %d: %s", name, e.Response.StatusCode, e.Response.Body)
				w.Header().Set("Content-Type", e.Response.ContentType)
				w.WriteHeader(e.Response.StatusCode)
				_, _ = w.Write([]byte(e.Response.Body))
				e.called++
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
	}))

	return ts
}

// AssertMockEndpointsCalls asserts that the number of calls an endpoint got
func AssertMockEndpointsCalls(endpoints []*MockEndpoint) error {
	for _, endpoint := range endpoints {
		if endpoint.called != endpoint.Calls {
			return fmt.Errorf(
				"expected endpoint %s %s to be called %d times but was called %d times",
				endpoint.Request.Method,
				endpoint.Request.Uri,
				endpoint.Calls,
				endpoint.called,
			)
		}
	}
	return nil
}

const (
	MockDataSourceWorkspacesBasic = `{
  "workspaces": [
    {
      "id": 1,
      "name": "Workspace 0001",
      "description": "Description",
      "apiKey": "691e0542-5c4d-4f74-be4a-38134a0aa0bf",
      "apiSecret": "8497f68e-75b9-431b-b067-cf86a074205c",
      "privateUrl": "/workspace/1",
      "publicUrl": "/share/1",
      "shareableUrl": ""
    }
  ]
}`
	MockResourceWorkspaceBasicCreate = `{
  "id": 1,
  "name": "Workspace 0001",
  "description": "Description",
  "apiKey": "691e0542-5c4d-4f74-be4a-38134a0aa0bf",
  "apiSecret": "8497f68e-75b9-431b-b067-cf86a074205c",
  "privateUrl": "/workspace/1",
  "publicUrl": "/share/1",
  "shareableUrl": ""
}`
	MockResourceWorkspaceBasicGet = `{
  "workspaces": [
    {
      "id": 1,
      "name": "Workspace 0001",
      "description": "Description",
      "apiKey": "691e0542-5c4d-4f74-be4a-38134a0aa0bf",
      "apiSecret": "8497f68e-75b9-431b-b067-cf86a074205c",
      "privateUrl": "/workspace/1",
      "publicUrl": "/share/1",
      "shareableUrl": ""
    }
  ]
}`
	MockResourceWorkspaceBasicDelete = `{
  "success": true,
  "message": "resource deleted successfully",
  "revision": 1
}`
	MockResourceWorkspaceWithSourceUpdate = `{
  "success": true,
  "message": "OK",
  "revision": 1
}`
	MockResourceWorkspaceWithSourceGet = `{
  "workspaces": [
    {
      "id": 1,
      "name": "Workspace DSL",
      "description": "Managed Workspace by DSL",
      "apiKey": "691e0542-5c4d-4f74-be4a-38134a0aa0bf",
      "apiSecret": "8497f68e-75b9-431b-b067-cf86a074205c",
      "privateUrl": "/workspace/1",
      "publicUrl": "/share/1",
      "shareableUrl": ""
    }
  ]
}`
)
