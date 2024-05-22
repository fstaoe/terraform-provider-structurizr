package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/api/model"
	"github.com/fstaoe/terraform-provider-structurizr/version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	DefaultUserAgent                 = "go-structurizr/" + version.LibraryVersion
	workspaceListCreateTemplate      = "/api/workspace"
	workspaceGetUpdateDeleteTemplate = "/api/workspace/%s"
)

// Config is the primary means to modify the Client
type Config struct {
	AdminAPIKey string
	BaseURL     *url.URL
	TLSInsecure bool
	UserAgent   string
}

// Client is the main Client API interface.
// Use NewClient to get started
type Client struct {
	config *Config
	doer   Doer
}

// Doer an interface that enables a more flexible dependency injection
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient returns an API Client used to communicate with the remote server
func NewClient(config *Config) *Client {
	return &Client{config, newHTTPClient(config.TLSInsecure)}
}

// GetWorkspaces lists all workspaces
func (c *Client) GetWorkspaces(ctx context.Context) (*model.Workspaces, error) {
	res, err := c.doCrud(
		ctx,
		http.MethodGet,
		workspaceListCreateTemplate,
		nil,
		new(model.Workspaces),
	)
	return res.(*model.Workspaces), err
}

// CreateWorkspace creates a new workspace
func (c *Client) CreateWorkspace(ctx context.Context) (*model.Workspace, error) {
	res, err := c.doCrud(
		ctx,
		http.MethodPost,
		workspaceListCreateTemplate,
		nil,
		new(model.Workspace),
	)
	return res.(*model.Workspace), err
}

// DeleteWorkspace deletes a workspace
func (c *Client) DeleteWorkspace(ctx context.Context, id int64) (*model.APIResponse, error) {
	u := urlEncodeTemplate(workspaceGetUpdateDeleteTemplate, strconv.FormatInt(id, 10))
	res, err := c.doCrud(ctx, http.MethodDelete, u, nil, new(model.APIResponse))
	return res.(*model.APIResponse), err
}

// newHTTPClient return an HTTP client configure TLS configuration for high customisation
func newHTTPClient(insecureSkipVerify bool) *http.Client {
	// Prevent issues with multiple data source configurations modifying the shared transport.
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureSkipVerify}

	return &http.Client{Transport: tr}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.config.BaseURL.ResolveReference(rel)
	var buf = new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("raw body to be sent over wire: %s", buf.String()))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Authorization", c.config.AdminAPIKey)
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	tflog.Debug(ctx, fmt.Sprintf("request: %+v", req))

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}

	// Drain and close the body to let the Transport reuse the connection
	// See https://github.com/google/go-github/pull/317 for more info/background
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	tflog.Trace(ctx, fmt.Sprintf("request: %+v response: %+v", req, resp))

	if resp.StatusCode >= 500 {
		return handleError(ctx, model.APIErrSystemUnavailable, req, resp)
	}

	if resp.StatusCode == 401 {
		return handleError(ctx, model.APIErrUnauthorized, req, resp)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return handleError(ctx, model.APIErrBadRequest, req, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to read the response body with: '%s'", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Received API response: %s", body))

	if v != nil {
		err = json.NewDecoder(io.NopCloser(bytes.NewBuffer(body))).Decode(v)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("failed decoding response for %s with: %s", req.URL.Path, err))
			return resp, err
		}
	}

	return resp, err
}

func (c *Client) doCrud(
	ctx context.Context,
	method string,
	path string,
	requestEntity interface{},
	responseEntity interface{},
) (interface{}, error) {
	var resp *http.Response

	req, err := c.newRequest(ctx, method, path, requestEntity)
	if err != nil {
		return responseEntity, err
	}

	if responseEntity == nil {
		if resp, err = c.do(ctx, req, nil); err == nil {
			// 201 -> extract the location header if the expectation is a string value
			switch resp.StatusCode {
			case http.StatusCreated:
				tflog.Trace(ctx, fmt.Sprintf("have 201, returning Location header: %+v", resp.Header))
				return resp.Header.Get("Location"), err
			}
		}
	} else {
		_, err = c.do(ctx, req, &responseEntity)
	}

	return responseEntity, err
}

func handleError(ctx context.Context, err error, req *http.Request, resp *http.Response) (*http.Response, error) {
	bodyBytes, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close() //  must close
	tflog.Debug(ctx, fmt.Sprintf("handling error response: %s", string(bodyBytes)))

	e := &model.APIErrorResponse{Err: err}
	decodingErr := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(e)
	if decodingErr != nil {
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"error decoding APIErrorResponse from response for %s %s with: %s",
				req.Method,
				req.URL.Path,
				decodingErr,
			),
		)
	}

	return resp, e
}

func urlEncodeTemplate(template string, parameters ...string) string {
	encodedParams := make([]interface{}, len(parameters))

	for i, p := range parameters {
		encodedParams[i] = url.PathEscape(p)
	}

	return fmt.Sprintf(template, encodedParams...)
}
