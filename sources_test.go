package volley_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/volleyhq/volley-go"
)

func TestListSources(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects/1/sources" {
			t.Errorf("Expected path /api/projects/1/sources, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"sources": []map[string]interface{}{
				{
					"id":                1,
					"slug":              "stripe-webhooks",
					"ingestion_id":      "src_abc123",
					"type":              "webhook",
					"eps":               10,
					"status":            "active",
					"connection_count":  2,
					"auth_type":         "none",
					"verify_signature":  false,
					"webhook_secret_set": false,
					"created_at":        time.Now().Format(time.RFC3339),
					"updated_at":        time.Now().Format(time.RFC3339),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	sources, err := client.ListSources(1)
	if err != nil {
		t.Fatalf("ListSources failed: %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("Expected 1 source, got %d", len(sources))
	}

	if sources[0].Slug != "stripe-webhooks" {
		t.Errorf("Expected source slug 'stripe-webhooks', got %s", sources[0].Slug)
	}
}

func TestCreateSource(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req volley.CreateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "Stripe Webhooks" {
			t.Errorf("Expected name 'Stripe Webhooks', got %s", req.Name)
		}

		if req.EPS != 10 {
			t.Errorf("Expected EPS 10, got %d", req.EPS)
		}

		response := map[string]interface{}{
			"source": map[string]interface{}{
				"id":            2,
				"slug":          "stripe-webhooks",
				"ingestion_id":  "src_xyz789",
				"type":          "webhook",
				"eps":           req.EPS,
				"status":        "active",
				"auth_type":     req.AuthType,
				"created_at":    time.Now().Format(time.RFC3339),
				"updated_at":    time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	source, err := client.CreateSource(1, volley.CreateSourceRequest{
		Name:     "Stripe Webhooks",
		EPS:      10,
		AuthType: "none",
	})
	if err != nil {
		t.Fatalf("CreateSource failed: %v", err)
	}

	if source.Slug != "stripe-webhooks" {
		t.Errorf("Expected source slug 'stripe-webhooks', got %s", source.Slug)
	}
}

func TestGetSource(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sources/1" {
			t.Errorf("Expected path /api/sources/1, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"source": map[string]interface{}{
				"id":           1,
				"slug":         "stripe-webhooks",
				"ingestion_id": "src_abc123",
				"eps":          10,
				"status":       "active",
				"created_at":   time.Now().Format(time.RFC3339),
				"updated_at":   time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	source, err := client.GetSource(1)
	if err != nil {
		t.Fatalf("GetSource failed: %v", err)
	}

	if source.ID != 1 {
		t.Errorf("Expected source ID 1, got %d", source.ID)
	}
}

func TestUpdateSource(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		var req volley.UpdateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		response := map[string]interface{}{
			"source": map[string]interface{}{
				"id":      1,
				"name":    req.Name,
				"eps":     *req.EPS,
				"status":  req.Status,
				"created_at": time.Now().Format(time.RFC3339),
				"updated_at": time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	eps := 20
	source, err := client.UpdateSource(1, volley.UpdateSourceRequest{
		Name:   "Updated Source",
		EPS:    &eps,
		Status: "active",
	})
	if err != nil {
		t.Fatalf("UpdateSource failed: %v", err)
	}

	if source.EPS != 20 {
		t.Errorf("Expected EPS 20, got %d", source.EPS)
	}
}

func TestDeleteSource(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/sources/1" {
			t.Errorf("Expected path /api/sources/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	err := client.DeleteSource(1)
	if err != nil {
		t.Fatalf("DeleteSource failed: %v", err)
	}
}

