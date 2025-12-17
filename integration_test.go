package volley_test

import (
	"os"
	"testing"

	"github.com/volleyhq/volley-go"
)

// Integration tests require VOLLEY_API_TOKEN environment variable
// These tests make real API calls to the Volley API
// Run with: go test -v -tags=integration -run TestIntegration

func TestIntegrationListOrganizations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiToken := os.Getenv("VOLLEY_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping integration test: VOLLEY_API_TOKEN not set")
	}

	client := volley.NewClient(apiToken)

	orgs, err := client.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}

	if len(orgs) == 0 {
		t.Log("No organizations found (this is OK if account is new)")
		return
	}

	t.Logf("Found %d organization(s)", len(orgs))
	for _, org := range orgs {
		t.Logf("  - %s (ID: %d, Role: %s)", org.Name, org.ID, org.Role)
	}
}

func TestIntegrationGetOrganization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiToken := os.Getenv("VOLLEY_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping integration test: VOLLEY_API_TOKEN not set")
	}

	client := volley.NewClient(apiToken)

	// First, list organizations to get an ID
	orgs, err := client.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}

	if len(orgs) == 0 {
		t.Skip("Skipping: No organizations found")
	}

	orgID := orgs[0].ID
	org, err := client.GetOrganization(&orgID)
	if err != nil {
		t.Fatalf("GetOrganization failed: %v", err)
	}

	if org.ID != orgID {
		t.Errorf("Expected organization ID %d, got %d", orgID, org.ID)
	}

	t.Logf("Retrieved organization: %s (ID: %d)", org.Name, org.ID)
}

func TestIntegrationListProjects(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiToken := os.Getenv("VOLLEY_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping integration test: VOLLEY_API_TOKEN not set")
	}

	client := volley.NewClient(apiToken)

	// Set organization context
	orgs, err := client.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}

	if len(orgs) == 0 {
		t.Skip("Skipping: No organizations found")
	}

	client.SetOrganizationID(orgs[0].ID)

	projects, err := client.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	t.Logf("Found %d project(s) in organization %s", len(projects), orgs[0].Name)
	for _, project := range projects {
		t.Logf("  - %s (ID: %d)", project.Name, project.ID)
	}
}

func TestIntegrationListSources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiToken := os.Getenv("VOLLEY_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping integration test: VOLLEY_API_TOKEN not set")
	}

	client := volley.NewClient(apiToken)

	// Set organization context and get a project
	orgs, err := client.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}

	if len(orgs) == 0 {
		t.Skip("Skipping: No organizations found")
	}

	client.SetOrganizationID(orgs[0].ID)

	projects, err := client.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects) == 0 {
		t.Skip("Skipping: No projects found")
	}

	sources, err := client.ListSources(projects[0].ID)
	if err != nil {
		t.Fatalf("ListSources failed: %v", err)
	}

	t.Logf("Found %d source(s) in project %s", len(sources), projects[0].Name)
	for _, source := range sources {
		t.Logf("  - %s (ID: %d, Ingestion ID: %s)", source.Slug, source.ID, source.IngestionID)
	}
}

