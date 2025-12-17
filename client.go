package volley

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default Volley API base URL
	DefaultBaseURL = "https://api.volleyhooks.com"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the main Volley API client
type Client struct {
	baseURL        string
	apiToken       string
	organizationID *uint64
	httpClient     *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new Volley API client
func NewClient(apiToken string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL:    DefaultBaseURL,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithOrganizationID sets the organization ID for all requests
func WithOrganizationID(orgID uint64) ClientOption {
	return func(c *Client) {
		c.organizationID = &orgID
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// SetOrganizationID sets the organization ID for subsequent requests
func (c *Client) SetOrganizationID(orgID uint64) {
	c.organizationID = &orgID
}

// ClearOrganizationID clears the organization ID
func (c *Client) ClearOrganizationID() {
	c.organizationID = nil
}

// BaseURL returns the base URL of the client (for testing)
func (c *Client) BaseURL() string {
	return c.baseURL
}

// OrganizationID returns the organization ID (for testing)
func (c *Client) OrganizationID() *uint64 {
	return c.organizationID
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}, queryParams map[string]string) (*http.Response, error) {
	// Build URL
	reqURL := c.baseURL + path
	if len(queryParams) > 0 {
		u, err := url.Parse(reqURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	// Create request body
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if c.organizationID != nil {
		req.Header.Set("X-Organization-ID", fmt.Sprintf("%d", *c.organizationID))
	}

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseResponse parses the JSON response into the target struct
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiError APIError
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return &apiError
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// APIError represents an API error response
type APIError struct {
	ErrorMsg string `json:"error"`
	Message  string `json:"message,omitempty"`
	Status   int    `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error: %s - %s", e.ErrorMsg, e.Message)
	}
	return fmt.Sprintf("API error: %s", e.ErrorMsg)
}

