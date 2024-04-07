package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/auth"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/model"
	"github.com/fstaoe/terraform-provider-structurizr/version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultUserAgent                 = "go-structurizr/" + version.LibraryVersion
	workspaceListCreateTemplate      = "/api/workspace"
	workspaceGetUpdateDeleteTemplate = "/api/workspace/%s"
)

// Config is the primary means to modify the Structurizr client
type Config struct {
	AdminAPIKey     string
	BaseURL         *url.URL
	CustomTLSConfig *tls.Config
	UserAgent       string
}

// Structurizr is the main Structurizr API interface.
// Use NewStructurizr to get started
type Structurizr struct {
	config     Config
	client     *http.Client
	timer      func() time.Time
	authFilter func(*http.Request)
}

// NewStructurizr creates a new Structurizr API client with sensible but overridable defaults
func NewStructurizr(config Config, client *http.Client, timer func() time.Time) *Structurizr {
	if client == nil {
		client = http.DefaultClient
	}

	if timer == nil {
		timer = func() time.Time { return time.Now() }
	}

	if config.CustomTLSConfig != nil && config.CustomTLSConfig.InsecureSkipVerify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Structurizr{config: config, client: client, timer: timer}
}

// WithAdminAuth adds an Authorization header using the Admin API key
func (c *Structurizr) WithAdminAuth() *Structurizr {
	c.authFilter = func(req *http.Request) {
		req.Header.Add("X-Authorization", c.config.AdminAPIKey)
	}
	return c
}

// WithHmacAuth adds an Authorization header using the HMAC scheme (https://en.wikipedia.org/wiki/HMAC)
func (c *Structurizr) WithHmacAuth(key string, secret string) *Structurizr {
	c.authFilter = func(req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			tflog.Error(req.Context(), fmt.Sprintf("failed to read the request body with: '%s'", err))
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		nonce := strconv.FormatInt(c.timer().UnixMilli(), 10)
		checksum := auth.Checksum(string(body))
		message := auth.Concat(
			req.Method,
			req.URL.Path,
			checksum,
			req.Header.Get("Content-Type"),
			nonce,
		)
		if len(body) > 0 {
			req.Header.Add("Content-MD5", auth.Enc(checksum))
		}
		req.Header.Add("X-Authorization", fmt.Sprintf("%s:%s", key, auth.Enc(auth.Sign(secret, message))))
		req.Header.Add("Nonce", nonce)
	}
	return c
}

// GetWorkspaces lists all workspaces
func (c *Structurizr) GetWorkspaces(ctx context.Context) (*model.Workspaces, error) {
	res, err := c.doCrud(ctx, http.MethodGet, workspaceListCreateTemplate, nil, new(model.Workspaces))
	return res.(*model.Workspaces), err
}

// CreateWorkspace creates a new workspace
func (c *Structurizr) CreateWorkspace(ctx context.Context) (*model.Workspace, error) {
	res, err := c.doCrud(ctx, http.MethodPost, workspaceListCreateTemplate, nil, new(model.Workspace))
	return res.(*model.Workspace), err
}

// DeleteWorkspace locks a workspace
func (c *Structurizr) DeleteWorkspace(ctx context.Context, id int64) (*model.GenericResponse, error) {
	u := urlEncodeTemplate(workspaceGetUpdateDeleteTemplate, strconv.FormatInt(id, 10))
	res, err := c.doCrud(ctx, http.MethodDelete, u, nil, new(model.GenericResponse))
	return res.(*model.GenericResponse), err
}

func (c *Structurizr) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.config.BaseURL.ResolveReference(rel)
	var buf = new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}

		tflog.Info(ctx, fmt.Sprintf("raw body to be sent over wire: '%s'", buf.String()))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	// Apply auth filter to the request
	c.authFilter(req)

	tflog.Debug(ctx, fmt.Sprintf("creating new request: %+v", req))
	return req, nil
}

func (c *Structurizr) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	tflog.Debug(ctx, fmt.Sprintf("sending body for request: %+v", req))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Drain and close the body to let the Transport reuse the connection
	// See https://github.com/google/go-github/pull/317/files for more info/background
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	tflog.Debug(ctx, fmt.Sprintf("response for request: %+v resp: %+v", req, resp))

	if resp.StatusCode >= 500 {
		return handleError(ctx, ErrSystemUnavailable, req, resp)
	}

	if resp.StatusCode == 401 {
		return handleError(ctx, ErrUnauthorized, req, resp)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return handleError(ctx, ErrBadRequest, req, resp)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("error decoding response for '%s'. Error: %s", req.URL.Path, err))
			return resp, err
		}
		tflog.Debug(ctx, fmt.Sprintf("response body: %+v", v))
	}

	return resp, err
}

func (c *Structurizr) doCrud(
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
				tflog.Debug(ctx, fmt.Sprintf("have 201, returning Location header: %+v", resp.Header))
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
	tflog.Debug(ctx, fmt.Sprintf("handling error response: '%s'", string(bodyBytes)))

	e := &apiErrorResponse{
		err: err,
	}
	decodingErr := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(e)
	if decodingErr != nil {
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"error decoding APIErrorResponse from response for '%s %s'. Error: %s",
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
