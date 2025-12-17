package volley

import (
	"fmt"
	"strconv"
	"time"
)

// ListDeliveryAttemptsOptions represents options for listing delivery attempts
type ListDeliveryAttemptsOptions struct {
	EventID       string
	SourceID      *uint64
	DestinationID *uint64
	ConnectionID  *uint64
	Status        string // "success" or "failed"
	StartTime     *time.Time
	EndTime       *time.Time
	Sort          string // "time", "time_oldest", "duration", "status_code"
	Limit         *int
	Offset        *int
}

// ListDeliveryAttempts lists all delivery attempts for a project with optional filters
func (c *Client) ListDeliveryAttempts(projectID uint64, opts *ListDeliveryAttemptsOptions) (*ListDeliveryAttemptsResponse, error) {
	path := fmt.Sprintf("/api/projects/%d/delivery-attempts", projectID)

	params := make(map[string]string)
	if opts != nil {
		if opts.EventID != "" {
			params["event_id"] = opts.EventID
		}
		if opts.SourceID != nil {
			params["source_id"] = strconv.FormatUint(*opts.SourceID, 10)
		}
		if opts.DestinationID != nil {
			params["destination_id"] = strconv.FormatUint(*opts.DestinationID, 10)
		}
		if opts.ConnectionID != nil {
			params["connection_id"] = strconv.FormatUint(*opts.ConnectionID, 10)
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
		if opts.Sort != "" {
			params["sort"] = opts.Sort
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

	var result ListDeliveryAttemptsResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

