package volley

// ListOrganizations lists all organizations the user has access to
func (c *Client) ListOrganizations() ([]Organization, error) {
	resp, err := c.doRequest("GET", "/api/org/list", nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Organizations []Organization `json:"organizations"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Organizations, nil
}

// GetOrganization gets the current organization
// If organizationID is provided, it will be used; otherwise, the first accessible organization is returned
func (c *Client) GetOrganization(organizationID *uint64) (*Organization, error) {
	// Temporarily set organization ID if provided
	originalOrgID := c.organizationID
	if organizationID != nil {
		c.organizationID = organizationID
	}

	resp, err := c.doRequest("GET", "/api/org", nil, nil)
	if err != nil {
		c.organizationID = originalOrgID
		return nil, err
	}

	var org Organization
	if err := c.parseResponse(resp, &org); err != nil {
		c.organizationID = originalOrgID
		return nil, err
	}

	c.organizationID = originalOrgID
	return &org, nil
}

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(req CreateOrganizationRequest) (*Organization, error) {
	resp, err := c.doRequest("POST", "/api/org", req, nil)
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := c.parseResponse(resp, &org); err != nil {
		return nil, err
	}

	return &org, nil
}

