package volley

import "time"

// Organization represents a Volley organization
type Organization struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	AccountID uint64    `json:"account_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Project represents a Volley project
type Project struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uint64    `json:"organization_id"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Source represents a webhook source
type Source struct {
	ID               uint64    `json:"id"`
	Slug             string    `json:"slug"`
	IngestionID      string    `json:"ingestion_id"`
	Type             string    `json:"type"`
	EPS              int       `json:"eps"`
	Status           string    `json:"status"`
	ConnectionCount  int64     `json:"connection_count"`
	AuthType         string    `json:"auth_type"`
	VerifySignature  bool      `json:"verify_signature"`
	WebhookSecretSet bool      `json:"webhook_secret_set"`
	AuthUsername     string    `json:"auth_username,omitempty"`
	AuthKeyName      string    `json:"auth_key_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Destination represents a webhook destination
type Destination struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	EPS       int       `json:"eps"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Connection represents a connection between a source and destination
type Connection struct {
	ID          uint64    `json:"id"`
	SourceID    uint64    `json:"source_id"`
	DestinationID uint64  `json:"destination_id"`
	Status      string    `json:"status"`
	EPS         int       `json:"eps"`
	MaxRetries  int       `json:"max_retries"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Event represents a webhook event/request
type Event struct {
	ID              uint64                 `json:"id"`
	EventID         string                 `json:"event_id"`
	SourceID        uint64                 `json:"source_id"`
	ProjectID       uint64                 `json:"project_id"`
	RawBody         string                 `json:"raw_body"`
	Headers         map[string]interface{} `json:"headers"`
	Status          string                 `json:"status"`
	DeliveryAttempts []DeliveryAttempt     `json:"delivery_attempts,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// DeliveryAttempt represents a delivery attempt for an event
type DeliveryAttempt struct {
	ID           uint64    `json:"id"`
	EventID      string    `json:"event_id"`
	ConnectionID uint64    `json:"connection_id"`
	Status       string    `json:"status"`
	StatusCode   int       `json:"status_code"`
	ErrorReason  string    `json:"error_reason,omitempty"`
	DurationMs   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// ListEventsResponse represents the response from listing events
type ListEventsResponse struct {
	PaginatedResponse
	Requests []Event `json:"requests"`
}

// ListDeliveryAttemptsResponse represents the response from listing delivery attempts
type ListDeliveryAttemptsResponse struct {
	PaginatedResponse
	Attempts []DeliveryAttempt `json:"attempts"`
}

// ReplayEventResponse represents the response from replaying an event
type ReplayEventResponse struct {
	Success    bool   `json:"success"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	ErrorReason string `json:"error_reason,omitempty"`
	DurationMs  int64  `json:"duration_ms"`
	AttemptID   uint64 `json:"attempt_id"`
}

