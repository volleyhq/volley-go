package volley

import "fmt"

// ListDestinations lists all destinations in a project
func (c *Client) ListDestinations(projectID uint64) ([]Destination, error) {
	path := fmt.Sprintf("/api/projects/%d/destinations", projectID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Destinations []Destination `json:"destinations"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Destinations, nil
}

// CreateDestinationRequest represents the request to create a destination
type CreateDestinationRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	EPS  int    `json:"eps"`
}

// CreateDestination creates a new destination
func (c *Client) CreateDestination(projectID uint64, req CreateDestinationRequest) (*Destination, error) {
	path := fmt.Sprintf("/api/projects/%d/destinations", projectID)
	resp, err := c.doRequest("POST", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Destination Destination `json:"destination"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Destination, nil
}

// GetDestination gets details of a specific destination
func (c *Client) GetDestination(destinationID uint64) (*Destination, error) {
	path := fmt.Sprintf("/api/destinations/%d", destinationID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Destination Destination `json:"destination"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Destination, nil
}

// UpdateDestinationRequest represents the request to update a destination
type UpdateDestinationRequest struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	EPS    *int   `json:"eps,omitempty"`
	Status string `json:"status,omitempty"`
}

// UpdateDestination updates a destination
func (c *Client) UpdateDestination(destinationID uint64, req UpdateDestinationRequest) (*Destination, error) {
	path := fmt.Sprintf("/api/destinations/%d", destinationID)
	resp, err := c.doRequest("PUT", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Destination Destination `json:"destination"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Destination, nil
}

// DeleteDestination deletes a destination
func (c *Client) DeleteDestination(destinationID uint64) error {
	path := fmt.Sprintf("/api/destinations/%d", destinationID)
	resp, err := c.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

