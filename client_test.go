package volley_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/volleyhooks/volley-go"
)

func TestNewClient(t *testing.T) {
	token := "test-token"
	client := volley.NewClient(token)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	// Test default values
	if client.BaseURL() != volley.DefaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", volley.DefaultBaseURL, client.BaseURL())
	}
}

func TestNewClientWithOptions(t *testing.T) {
	token := "test-token"
	customURL := "https://api-staging.volleyhooks.com"
	orgID := uint64(123)

	client := volley.NewClient(token,
		volley.WithBaseURL(customURL),
		volley.WithOrganizationID(orgID),
	)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.BaseURL() != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.BaseURL())
	}

	if client.OrganizationID() == nil || *client.OrganizationID() != orgID {
		t.Errorf("Expected organization ID %d, got %v", orgID, client.OrganizationID())
	}
}

func TestSetOrganizationID(t *testing.T) {
	client := volley.NewClient("test-token")

	orgID := uint64(456)
	client.SetOrganizationID(orgID)

	if client.OrganizationID() == nil || *client.OrganizationID() != orgID {
		t.Errorf("Expected organization ID %d, got %v", orgID, client.OrganizationID())
	}
}

func TestClearOrganizationID(t *testing.T) {
	client := volley.NewClient("test-token")
	client.SetOrganizationID(123)
	client.ClearOrganizationID()

	if client.OrganizationID() != nil {
		t.Error("Expected organization ID to be nil after ClearOrganizationID")
	}
}

func TestClientWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	client := volley.NewClient("test-token",
		volley.WithHTTPClient(customClient),
	)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

// Helper function to create a test server
func createTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

