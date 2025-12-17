package volley_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/volleyhooks/volley-go"
)

func TestListEvents(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects/1/requests" {
			t.Errorf("Expected path /api/projects/1/requests, got %s", r.URL.Path)
		}

		// Verify query parameters
		status := r.URL.Query().Get("status")
		if status != "failed" {
			t.Errorf("Expected status query param 'failed', got %s", status)
		}

		limit := r.URL.Query().Get("limit")
		if limit != "50" {
			t.Errorf("Expected limit query param '50', got %s", limit)
		}

		response := map[string]interface{}{
			"requests": []map[string]interface{}{
				{
					"id":        1,
					"event_id":  "evt_abc123",
					"source_id": 10,
					"status":    "failed",
					"created_at": time.Now().Format(time.RFC3339),
				},
			},
			"total":  1,
			"limit": 50,
			"offset": 0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	limit := 50
	events, err := client.ListEvents(1, &volley.ListEventsOptions{
		Status: "failed",
		Limit:  &limit,
	})
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}

	if events.Total != 1 {
		t.Errorf("Expected total 1, got %d", events.Total)
	}

	if len(events.Requests) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events.Requests))
	}

	if events.Requests[0].EventID != "evt_abc123" {
		t.Errorf("Expected event ID 'evt_abc123', got %s", events.Requests[0].EventID)
	}
}

func TestGetEvent(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/requests/123" {
			t.Errorf("Expected path /api/requests/123, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"request": map[string]interface{}{
				"id":        123,
				"event_id":  "evt_abc123",
				"source_id": 10,
				"project_id": 1,
				"raw_body":  `{"event": "user.created"}`,
				"headers":   map[string]interface{}{"Content-Type": "application/json"},
				"status":    "failed",
				"created_at": time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	event, err := client.GetEvent(123)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}

	if event.EventID != "evt_abc123" {
		t.Errorf("Expected event ID 'evt_abc123', got %s", event.EventID)
	}

	if event.Status != "failed" {
		t.Errorf("Expected status 'failed', got %s", event.Status)
	}
}

func TestReplayEvent(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/replay-event" {
			t.Errorf("Expected path /api/replay-event, got %s", r.URL.Path)
		}

		var req volley.ReplayEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.EventID != "evt_abc123" {
			t.Errorf("Expected event ID 'evt_abc123', got %s", req.EventID)
		}

		response := map[string]interface{}{
			"success":      true,
			"status":       "success",
			"status_code":  200,
			"duration_ms":  150,
			"attempt_id":   456,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	result, err := client.ReplayEvent(volley.ReplayEventRequest{
		EventID: "evt_abc123",
	})
	if err != nil {
		t.Fatalf("ReplayEvent failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected replay to be successful")
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
}

func TestListEventsWithFilters(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		// Verify all query parameters
		sourceID := r.URL.Query().Get("source_id")
		if sourceID != "10" {
			t.Errorf("Expected source_id '10', got %s", sourceID)
		}

		startTime := r.URL.Query().Get("start_time")
		if startTime == "" {
			t.Error("Expected start_time query parameter")
		}

		search := r.URL.Query().Get("search")
		if search != "test" {
			t.Errorf("Expected search 'test', got %s", search)
		}

		response := map[string]interface{}{
			"requests": []map[string]interface{}{},
			"total":    0,
			"limit":    50,
			"offset":   0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	sourceID := uint64(10)
	startTime := time.Now().Add(-24 * time.Hour)
	limit := 50

	_, err := client.ListEvents(1, &volley.ListEventsOptions{
		SourceID:  &sourceID,
		StartTime: &startTime,
		Search:    "test",
		Limit:     &limit,
	})
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
}

