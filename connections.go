package volley

import "fmt"

// CreateConnectionRequest represents the request to create a connection
type CreateConnectionRequest struct {
	SourceID    uint64 `json:"source_id"`
	DestinationID uint64 `json:"destination_id"`
	Status      string `json:"status"` // "enabled" or "disabled"
	EPS         int    `json:"eps"`
	MaxRetries  int    `json:"max_retries"`
}

// CreateConnection creates a connection between a source and destination
func (c *Client) CreateConnection(projectID uint64, req CreateConnectionRequest) (*Connection, error) {
	path := fmt.Sprintf("/api/projects/%d/connections", projectID)
	resp, err := c.doRequest("POST", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Connection Connection `json:"connection"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Connection, nil
}

// GetConnection gets details and metrics for a connection
func (c *Client) GetConnection(connectionID uint64) (*Connection, error) {
	path := fmt.Sprintf("/api/connections/%d", connectionID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Connection Connection `json:"connection"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Connection, nil
}

// UpdateConnectionRequest represents the request to update a connection
type UpdateConnectionRequest struct {
	Status     string `json:"status,omitempty"` // "enabled" or "disabled"
	EPS        *int   `json:"eps,omitempty"`
	MaxRetries *int   `json:"max_retries,omitempty"`
}

// UpdateConnection updates a connection
func (c *Client) UpdateConnection(connectionID uint64, req UpdateConnectionRequest) (*Connection, error) {
	path := fmt.Sprintf("/api/connections/%d", connectionID)
	resp, err := c.doRequest("PUT", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Connection Connection `json:"connection"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Connection, nil
}

// DeleteConnection deletes a connection
func (c *Client) DeleteConnection(connectionID uint64) error {
	path := fmt.Sprintf("/api/connections/%d", connectionID)
	resp, err := c.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

