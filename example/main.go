package main

import (
	"fmt"
	"log"
	"os"

	"github.com/volleyhq/volley-go"
)

func main() {
	// Get API token from environment variable
	apiToken := os.Getenv("VOLLEY_API_TOKEN")
	if apiToken == "" {
		log.Fatal("VOLLEY_API_TOKEN environment variable is required")
	}

	// Create a client
	client := volley.NewClient(apiToken)

	// Example 1: List organizations
	fmt.Println("=== Listing Organizations ===")
	orgs, err := client.ListOrganizations()
	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	for _, org := range orgs {
		fmt.Printf("  - %s (ID: %d, Role: %s)\n", org.Name, org.ID, org.Role)
	}

	if len(orgs) == 0 {
		fmt.Println("  No organizations found")
		return
	}

	// Example 2: Set organization context
	orgID := orgs[0].ID
	client.SetOrganizationID(orgID)
	fmt.Printf("\n=== Using Organization: %s (ID: %d) ===\n", orgs[0].Name, orgID)

	// Example 3: List projects
	fmt.Println("\n=== Listing Projects ===")
	projects, err := client.ListProjects()
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	for _, project := range projects {
		fmt.Printf("  - %s (ID: %d", project.Name, project.ID)
		if project.IsDefault {
			fmt.Print(", Default")
		}
		fmt.Println(")")
	}

	if len(projects) == 0 {
		fmt.Println("  No projects found")
		return
	}

	// Example 4: List sources for first project
	projectID := projects[0].ID
	fmt.Printf("\n=== Listing Sources for Project: %s (ID: %d) ===\n", projects[0].Name, projectID)
	sources, err := client.ListSources(projectID)
	if err != nil {
		log.Fatalf("Failed to list sources: %v", err)
	}

	for _, source := range sources {
		fmt.Printf("  - %s (ID: %d, Ingestion ID: %s, Type: %s)\n",
			source.Slug, source.ID, source.IngestionID, source.Type)
	}

	// Example 5: List events (if any)
	fmt.Printf("\n=== Listing Recent Events for Project: %s ===\n", projects[0].Name)
	limit := 10
	events, err := client.ListEvents(projectID, &volley.ListEventsOptions{
		Limit: &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list events: %v", err)
	}

	fmt.Printf("Total events: %d\n", events.Total)
	for i, event := range events.Requests {
		if i >= 5 { // Show only first 5
			break
		}
		fmt.Printf("  - Event ID: %s, Status: %s\n", event.EventID, event.Status)
	}
}

