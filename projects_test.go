package volley_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/volleyhooks/volley-go"
)

func TestListProjects(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects" {
			t.Errorf("Expected path /api/projects, got %s", r.URL.Path)
		}

		// Verify X-Organization-ID header
		orgIDHeader := r.Header.Get("X-Organization-ID")
		if orgIDHeader != "123" {
			t.Errorf("Expected X-Organization-ID header '123', got %s", orgIDHeader)
		}

		response := map[string]interface{}{
			"projects": []map[string]interface{}{
				{
					"id":              1,
					"name":            "Test Project",
					"organization_id": 123,
					"is_default":      true,
					"created_at":      time.Now().Format(time.RFC3339),
					"updated_at":      time.Now().Format(time.RFC3339),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))
	client.SetOrganizationID(123)

	projects, err := client.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(projects))
	}

	if projects[0].Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", projects[0].Name)
	}
}

func TestCreateProject(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req volley.CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "New Project" {
			t.Errorf("Expected name 'New Project', got %s", req.Name)
		}

		response := map[string]interface{}{
			"project": map[string]interface{}{
				"id":              2,
				"name":            req.Name,
				"organization_id": 123,
				"is_default":      req.IsDefault,
				"created_at":      time.Now().Format(time.RFC3339),
				"updated_at":      time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	project, err := client.CreateProject(volley.CreateProjectRequest{
		Name:      "New Project",
		IsDefault: false,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if project.Name != "New Project" {
		t.Errorf("Expected project name 'New Project', got %s", project.Name)
	}
}

func TestUpdateProject(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects/1" {
			t.Errorf("Expected path /api/projects/1, got %s", r.URL.Path)
		}

		var req volley.UpdateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		response := map[string]interface{}{
			"project": map[string]interface{}{
				"id":              1,
				"name":            req.Name,
				"organization_id": 123,
				"is_default":      false,
				"created_at":      time.Now().Format(time.RFC3339),
				"updated_at":      time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	project, err := client.UpdateProject(1, volley.UpdateProjectRequest{
		Name: "Updated Project Name",
	})
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	if project.Name != "Updated Project Name" {
		t.Errorf("Expected project name 'Updated Project Name', got %s", project.Name)
	}
}

func TestDeleteProject(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects/1" {
			t.Errorf("Expected path /api/projects/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	err := client.DeleteProject(1)
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
}

func TestGetConnections(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/projects/1/connections" {
			t.Errorf("Expected path /api/projects/1/connections, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"connections": []map[string]interface{}{
				{
					"id":            1,
					"source_id":     10,
					"destination_id": 20,
					"status":        "enabled",
					"eps":           5,
					"max_retries":   3,
					"created_at":    time.Now().Format(time.RFC3339),
					"updated_at":    time.Now().Format(time.RFC3339),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))

	connections, err := client.GetConnections(1)
	if err != nil {
		t.Fatalf("GetConnections failed: %v", err)
	}

	if len(connections) != 1 {
		t.Fatalf("Expected 1 connection, got %d", len(connections))
	}

	if connections[0].Status != "enabled" {
		t.Errorf("Expected connection status 'enabled', got %s", connections[0].Status)
	}
}

