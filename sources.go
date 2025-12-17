package volley

import "fmt"

// ListSources lists all sources in a project
func (c *Client) ListSources(projectID uint64) ([]Source, error) {
	path := fmt.Sprintf("/api/projects/%d/sources", projectID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Sources []Source `json:"sources"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Sources, nil
}

// CreateSourceRequest represents the request to create a source
type CreateSourceRequest struct {
	Name     string `json:"name"`
	EPS      int    `json:"eps"`
	AuthType string `json:"auth_type"` // "none", "basic", "api_key"
}

// CreateSource creates a new source
func (c *Client) CreateSource(projectID uint64, req CreateSourceRequest) (*Source, error) {
	path := fmt.Sprintf("/api/projects/%d/sources", projectID)
	resp, err := c.doRequest("POST", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Source Source `json:"source"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Source, nil
}

// GetSource gets details of a specific source
func (c *Client) GetSource(sourceID uint64) (*Source, error) {
	path := fmt.Sprintf("/api/sources/%d", sourceID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Source Source `json:"source"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Source, nil
}

// UpdateSourceRequest represents the request to update a source
type UpdateSourceRequest struct {
	Name     string `json:"name,omitempty"`
	EPS      *int   `json:"eps,omitempty"`
	AuthType string `json:"auth_type,omitempty"`
	Status   string `json:"status,omitempty"`
}

// UpdateSource updates a source
func (c *Client) UpdateSource(sourceID uint64, req UpdateSourceRequest) (*Source, error) {
	path := fmt.Sprintf("/api/sources/%d", sourceID)
	resp, err := c.doRequest("PUT", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Source Source `json:"source"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Source, nil
}

// DeleteSource deletes a source
func (c *Client) DeleteSource(sourceID uint64) error {
	path := fmt.Sprintf("/api/sources/%d", sourceID)
	resp, err := c.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

