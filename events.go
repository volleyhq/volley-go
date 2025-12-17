package volley

import (
	"fmt"
	"strconv"
	"time"
)

// ListEventsOptions represents options for listing events
type ListEventsOptions struct {
	SourceID      *uint64
	ConnectionID  *uint64
	DestinationID *uint64
	Status        string // "processed", "pending", "failed", "dropped"
	StartTime     *time.Time
	EndTime       *time.Time
	Search        string
	Limit         *int
	Offset        *int
}

// ListEvents lists all events/requests for a project with optional filters
func (c *Client) ListEvents(projectID uint64, opts *ListEventsOptions) (*ListEventsResponse, error) {
	path := fmt.Sprintf("/api/projects/%d/requests", projectID)

	params := make(map[string]string)
	if opts != nil {
		if opts.SourceID != nil {
			params["source_id"] = strconv.FormatUint(*opts.SourceID, 10)
		}
		if opts.ConnectionID != nil {
			params["connection_id"] = strconv.FormatUint(*opts.ConnectionID, 10)
		}
		if opts.DestinationID != nil {
			params["destination_id"] = strconv.FormatUint(*opts.DestinationID, 10)
		}
		if opts.Status != "" {
			params["status"] = opts.Status
		}
		if opts.StartTime != nil {
			params["start_time"] = opts.StartTime.Format(time.RFC3339)
		}
		if opts.EndTime != nil {
			params["end_time"] = opts.EndTime.Format(time.RFC3339)
		}
		if opts.Search != "" {
			params["search"] = opts.Search
		}
		if opts.Limit != nil {
			params["limit"] = strconv.Itoa(*opts.Limit)
		}
		if opts.Offset != nil {
			params["offset"] = strconv.Itoa(*opts.Offset)
		}
	}

	resp, err := c.doRequest("GET", path, nil, params)
	if err != nil {
		return nil, err
	}

	var result ListEventsResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetEvent gets detailed information about a specific event by its database ID
func (c *Client) GetEvent(requestID uint64) (*Event, error) {
	path := fmt.Sprintf("/api/requests/%d", requestID)
	resp, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Request Event `json:"request"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Request, nil
}

// ReplayEventRequest represents the request to replay an event
type ReplayEventRequest struct {
	EventID       string  `json:"event_id"`
	DestinationID *uint64 `json:"destination_id,omitempty"`
	ConnectionID  *uint64 `json:"connection_id,omitempty"`
}

// ReplayEvent replays a failed event by its event_id
func (c *Client) ReplayEvent(req ReplayEventRequest) (*ReplayEventResponse, error) {
	resp, err := c.doRequest("POST", "/api/replay-event", req, nil)
	if err != nil {
		return nil, err
	}

	var result ReplayEventResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

