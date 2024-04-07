package client

import (
	"bytes"
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewStructurizr(t *testing.T) {
	config := Config{CustomTLSConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{}
	timer := func() time.Time { return time.UnixMilli(1713472217559) }

	t.Run("Given all arguments", func(t *testing.T) {
		got := NewStructurizr(config, client, timer)
		assert.Equal(t, config, got.config)
		assert.Equal(t, client, got.client)
		assert.Equal(t, timer(), got.timer())
	})
	t.Run("Given TLS config", func(t *testing.T) {
		config.CustomTLSConfig.InsecureSkipVerify = false
		got := NewStructurizr(config, client, timer)
		assert.Equal(t, config, got.config)
		assert.Equal(t, client, got.client)
		assert.IsType(t, timer(), got.timer())
	})
	t.Run("Given default HTTP client", func(t *testing.T) {
		got := NewStructurizr(config, nil, timer)
		assert.Equal(t, config, got.config)
		assert.Equal(t, http.DefaultClient, got.client)
		assert.Equal(t, timer(), got.timer())
	})
	t.Run("Given default timer", func(t *testing.T) {
		got := NewStructurizr(config, client, nil)
		assert.Equal(t, config, got.config)
		assert.Equal(t, client, got.client)
		assert.IsType(t, time.Time{}, got.timer())
	})
}

func TestStructurizr_WithAPIKeyAuth(t *testing.T) {
	apiKey := "structurizr"
	client := &Structurizr{config: Config{AdminAPIKey: apiKey}}
	req := &http.Request{Header: make(http.Header)}
	client.WithAdminAuth().authFilter(req)
	assert.Equal(t, apiKey, req.Header.Get("X-Authorization"))
}

func TestStructurizr_WithHmacAuth(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080")
	type fields struct {
		config Config
		timer  func() time.Time
	}
	type args struct {
		apiKey    string
		apiSecret string
		req       *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			"Given a GET request by ID",
			fields{
				config: Config{
					AdminAPIKey:     "structurizr",
					BaseURL:         baseURL,
					CustomTLSConfig: &tls.Config{InsecureSkipVerify: true},
				},
				timer: func() time.Time { return time.UnixMilli(1713472217559) },
			},
			args{
				apiKey:    "b9ca9ee8-6917-4447-8654-39d8eba0447d",
				apiSecret: "f17ec048-e92d-48f9-a76e-9b2716da807c",
				req: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: baseURL.Scheme, Host: baseURL.Host, Path: "/api/workspace/1"},
					Header: make(http.Header),
					Body:   io.NopCloser(new(bytes.Buffer)),
				},
			},
			map[string]string{
				"X-Authorization": "b9ca9ee8-6917-4447-8654-39d8eba0447d:NGVkMjNmOWEyYjVkNDE3YjYwMTBiYjk2NDY0OWQ2ZmM0ODA4YmY2MDdjNDk2MDMwMjZkZTkzNmNhOGZhOWVjMA==",
				"Nonce":           "1713472217559",
			},
		},
		{
			"Given a PUT request",
			fields{
				config: Config{
					AdminAPIKey:     "structurizr",
					BaseURL:         baseURL,
					CustomTLSConfig: &tls.Config{InsecureSkipVerify: true},
				},
				timer: func() time.Time { return time.UnixMilli(1713472217559) },
			},
			args{
				apiKey:    "b9ca9ee8-6917-4447-8654-39d8eba0447d",
				apiSecret: "f17ec048-e92d-48f9-a76e-9b2716da807c",
				req: &http.Request{
					Method: http.MethodPut,
					URL:    &url.URL{Scheme: baseURL.Scheme, Host: baseURL.Host, Path: "/api/workspace/1"},
					Header: make(http.Header),
					Body:   io.NopCloser(bytes.NewBuffer([]byte("{\"id\": 1}"))),
				},
			},
			map[string]string{
				"Content-MD5":     "ZjNlNTZjNjAyNzcxZTk1NDFhZWY2MWQ1MDI1NjJiODk=",
				"X-Authorization": "b9ca9ee8-6917-4447-8654-39d8eba0447d:NDk2NjViZjBmN2ViNjQ5YzlkZTU5NGViYTIwZWU1ZTA5OTkyODkwNTEyYjE4MWRmZWQyZTQ3NDI2MTkzNDQ2Nw==",
				"Nonce":           "1713472217559",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Structurizr{config: tt.fields.config, timer: tt.fields.timer}
			// Attach HMAC scheme authentication to the client
			c.WithHmacAuth(tt.args.apiKey, tt.args.apiSecret)
			// Apply the auth filter to the request
			c.authFilter(tt.args.req)
			// Check
			for k, want := range tt.want {
				got := tt.args.req.Header.Get(k)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("WithHmacAuth() = %v, want %v", got, want)
				}
			}
		})
	}
}
