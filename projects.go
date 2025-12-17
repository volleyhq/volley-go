package volley

import "fmt"

// ListProjects lists all projects in the current organization
func (c *Client) ListProjects() ([]Project, error) {
	resp, err := c.doRequest("GET", "/api/projects", nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Projects []Project `json:"projects"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Projects, nil
}

// CreateProjectRequest represents the request to create a project
type CreateProjectRequest struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default,omitempty"`
}

// CreateProject creates a new project
func (c *Client) CreateProject(req CreateProjectRequest) (*Project, error) {
	resp, err := c.doRequest("POST", "/api/projects", req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Project Project `json:"project"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Project, nil
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateProject updates a project's name
func (c *Client) UpdateProject(projectID uint64, req UpdateProjectRequest) (*Project, error) {
	path := fmt.Sprintf("/api/projects/%d", projectID)
	resp, err := c.doRequest("PUT", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Project Project `json:"project"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID uint64) error {
	path := fmt.Sprintf("/api/projects/%d", projectID)
	resp, err := c.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetConnections lists all connections in a project
func (c *Client) GetConnections(projectID uint64) ([]Connection, error) {
	path := fmt.Sprintf("/api/projects/%d/connections", projectID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Connections []Connection `json:"connections"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Connections, nil
}

