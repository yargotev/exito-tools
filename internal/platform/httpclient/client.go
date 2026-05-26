package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the default timeout for provider HTTP clients.
	DefaultTimeout = 10 * time.Second
	// MaxResponseBodyBytes is the maximum provider response body size decoded by shared helpers.
	MaxResponseBodyBytes = 1 << 20
)

// Config contains provider-agnostic settings for shared HTTP clients.
type Config struct {
	BaseURL string
	Token   string
	Timeout time.Duration
	Client  *http.Client
}

// Client wraps low-level HTTP concerns shared by domain-owned provider clients.
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// New creates a shared HTTP client from provider-agnostic settings.
func New(config Config) Client {
	client := config.Client
	if client == nil {
		timeout := config.Timeout
		if timeout <= 0 {
			timeout = DefaultTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	return Client{
		baseURL: strings.TrimSpace(config.BaseURL),
		token:   strings.TrimSpace(config.Token),
		client:  client,
	}
}

// Configured reports whether the client has the minimum provider settings needed to call authenticated APIs.
func (c Client) Configured() bool {
	return c.baseURL != "" && c.token != ""
}

// NewJSONRequest builds an authenticated JSON provider request and applies execution metadata headers.
func (c Client) NewJSONRequest(ctx context.Context, method string, path string, payload any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := c.NewRequest(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// NewRequest builds an authenticated provider request and applies execution metadata headers.
func (c Client) NewRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	endpoint, err := c.Endpoint(path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	if c.token != "" {
		request.Header.Set("Authorization", "Bearer "+c.token)
	}
	ApplyRequestMetadata(ctx, request)

	return request, nil
}

// Endpoint resolves a provider-relative path against the configured base URL.
func (c Client) Endpoint(path string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(c.baseURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("absolute URL required")
	}

	parsed.Path = joinPath(parsed.Path, path)
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String(), nil
}

// Do sends a prepared HTTP request using the configured underlying client.
func (c Client) Do(request *http.Request) (*http.Response, error) {
	return c.client.Do(request) // #nosec G107,G704 -- domain clients pass requests resolved from configured provider base URLs.
}

// Successful reports whether an HTTP response has a 2xx status code.
func Successful(response *http.Response) bool {
	return response != nil && response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
}

// DecodeJSONResponse decodes a bounded JSON response body into target.
func DecodeJSONResponse(response *http.Response, target any) error {
	limited := io.LimitReader(response.Body, MaxResponseBodyBytes)
	return json.NewDecoder(limited).Decode(target)
}

func joinPath(basePath string, path string) string {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		if basePath == "" {
			return "/"
		}
		return basePath
	}

	return strings.TrimRight(basePath, "/") + "/" + strings.TrimLeft(trimmedPath, "/")
}
