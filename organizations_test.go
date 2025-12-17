package volley_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/volleyhq/volley-go"
)

func TestListOrganizations(t *testing.T) {
	// Create a test server
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/org/list" {
			t.Errorf("Expected path /api/org/list, got %s", r.URL.Path)
		}

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
		}

		response := map[string]interface{}{
			"organizations": []map[string]interface{}{
				{
					"id":         1,
					"name":       "Test Org",
					"slug":       "test-org",
					"account_id": 100,
					"role":       "owner",
					"created_at": time.Now().Format(time.RFC3339),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	// Create client with test server URL
	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	orgs, err := client.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}

	if len(orgs) != 1 {
		t.Fatalf("Expected 1 organization, got %d", len(orgs))
	}

	if orgs[0].Name != "Test Org" {
		t.Errorf("Expected organization name 'Test Org', got %s", orgs[0].Name)
	}

	if orgs[0].ID != 1 {
		t.Errorf("Expected organization ID 1, got %d", orgs[0].ID)
	}
}

func TestGetOrganization(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/org" {
			t.Errorf("Expected path /api/org, got %s", r.URL.Path)
		}

		// Verify X-Organization-ID header if set
		orgIDHeader := r.Header.Get("X-Organization-ID")
		if orgIDHeader != "" && orgIDHeader != "123" {
			t.Errorf("Expected X-Organization-ID header '123', got %s", orgIDHeader)
		}

		response := map[string]interface{}{
			"id":   123,
			"name": "Test Org",
			"slug": "test-org",
			"role": "owner",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	// Test with organization ID
	orgID := uint64(123)
	org, err := client.GetOrganization(&orgID)
	if err != nil {
		t.Fatalf("GetOrganization failed: %v", err)
	}

	if org.ID != 123 {
		t.Errorf("Expected organization ID 123, got %d", org.ID)
	}

	// Test without organization ID
	org, err = client.GetOrganization(nil)
	if err != nil {
		t.Fatalf("GetOrganization failed: %v", err)
	}

	if org.Name != "Test Org" {
		t.Errorf("Expected organization name 'Test Org', got %s", org.Name)
	}
}

func TestCreateOrganization(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/org" {
			t.Errorf("Expected path /api/org, got %s", r.URL.Path)
		}

		// Parse request body
		var req volley.CreateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "New Organization" {
			t.Errorf("Expected name 'New Organization', got %s", req.Name)
		}

		response := map[string]interface{}{
			"id":         2,
			"name":       req.Name,
			"slug":       "new-organization",
			"account_id": 100,
			"role":       "owner",
			"created_at": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	org, err := client.CreateOrganization(volley.CreateOrganizationRequest{
		Name: "New Organization",
	})
	if err != nil {
		t.Fatalf("CreateOrganization failed: %v", err)
	}

	if org.Name != "New Organization" {
		t.Errorf("Expected organization name 'New Organization', got %s", org.Name)
	}

	if org.ID != 2 {
		t.Errorf("Expected organization ID 2, got %d", org.ID)
	}
}

func TestListOrganizationsError(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "unauthorized",
		})
	})
	defer server.Close()

	client := volley.NewClient("invalid-token", volley.WithBaseURL(server.URL))

	_, err := client.ListOrganizations()
	if err == nil {
		t.Fatal("Expected error for unauthorized request")
	}

	apiErr, ok := err.(*volley.APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.ErrorMsg != "unauthorized" {
		t.Errorf("Expected error message 'unauthorized', got %s", apiErr.ErrorMsg)
	}
}

