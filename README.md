# Volley Go SDK

Official Go SDK for the Volley API. This SDK provides a convenient way to interact with the Volley webhook infrastructure API.

[Volley](https://volleyhooks.com) is a webhook infrastructure platform that provides reliable webhook delivery, rate limiting, retries, monitoring, and more.

## Resources

- **Documentation**: [https://docs.volleyhooks.com](https://docs.volleyhooks.com)
- **Getting Started Guide**: [https://docs.volleyhooks.com/getting-started](https://docs.volleyhooks.com/getting-started)
- **API Reference**: [https://docs.volleyhooks.com/api](https://docs.volleyhooks.com/api)
- **Authentication Guide**: [https://docs.volleyhooks.com/authentication](https://docs.volleyhooks.com/authentication)
- **Security Guide**: [https://docs.volleyhooks.com/security](https://docs.volleyhooks.com/security)
- **Console**: [https://app.volleyhooks.com](https://app.volleyhooks.com)
- **Website**: [https://volleyhooks.com](https://volleyhooks.com)

## Installation

```bash
go get github.com/volleyhooks/volley-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/volleyhooks/volley-go"
)

func main() {
    // Create a client with your API token
    client := volley.NewClient("your-api-token")
    
    // Optionally set organization context
    orgID := uint64(123)
    client.SetOrganizationID(orgID)
    
    // List organizations
    orgs, err := client.ListOrganizations()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, org := range orgs {
        fmt.Printf("Organization: %s (ID: %d)\n", org.Name, org.ID)
    }
}
```

## Authentication

Volley uses API tokens for authentication. These are long-lived tokens designed for programmatic access.

### Getting Your API Token

1. Log in to the [Volley Console](https://app.volleyhooks.com)
2. Navigate to **Settings → Account → API Token**
3. Click **View Token** (you may need to verify your password)
4. Copy the token and store it securely

**Important**: API tokens are non-expiring and provide full access to your account. Keep them secure and rotate them if compromised. See the [Security Guide](https://docs.volleyhooks.com/security) for best practices.

```go
client := volley.NewClient("your-api-token")
```

For more details on authentication, API tokens, and security, see the [Authentication Guide](https://docs.volleyhooks.com/authentication) and [Security Guide](https://docs.volleyhooks.com/security).

## Organization Context

When you have multiple organizations, you need to specify which organization context to use for API requests. The API verifies that resources (like projects) belong to the specified organization.

You can set the organization context in two ways:

```go
// Method 1: Set organization ID for all subsequent requests
client.SetOrganizationID(123)

// Method 2: Create client with organization ID
client := volley.NewClient("your-api-token", 
    volley.WithOrganizationID(123),
)

// Clear organization context (uses first accessible organization)
client.ClearOrganizationID()
```

**Note**: If you don't set an organization ID, the API uses your first accessible organization by default. For more details, see the [API Reference - Organization Context](https://docs.volleyhooks.com/api#organization-context).

## Examples

### Organizations

```go
// List all organizations
orgs, err := client.ListOrganizations()
if err != nil {
    log.Fatal(err)
}

// Get current organization
org, err := client.GetOrganization(nil) // nil = use default
if err != nil {
    log.Fatal(err)
}

// Create organization
newOrg, err := client.CreateOrganization(volley.CreateOrganizationRequest{
    Name: "My Organization",
})
if err != nil {
    log.Fatal(err)
}
```

### Projects

```go
// List projects
projects, err := client.ListProjects()
if err != nil {
    log.Fatal(err)
}

// Create project
project, err := client.CreateProject(volley.CreateProjectRequest{
    Name:      "My Project",
    IsDefault: false,
})
if err != nil {
    log.Fatal(err)
}

// Update project
updated, err := client.UpdateProject(project.ID, volley.UpdateProjectRequest{
    Name: "Updated Name",
})
if err != nil {
    log.Fatal(err)
}

// Delete project
err = client.DeleteProject(project.ID)
if err != nil {
    log.Fatal(err)
}
```

### Sources

```go
// List sources in a project
sources, err := client.ListSources(projectID)
if err != nil {
    log.Fatal(err)
}

// Create source
source, err := client.CreateSource(projectID, volley.CreateSourceRequest{
    Name:     "Stripe Webhooks",
    EPS:      10,
    AuthType: "none",
})
if err != nil {
    log.Fatal(err)
}

// Get source details
source, err := client.GetSource(sourceID)
if err != nil {
    log.Fatal(err)
}

// Update source
updated, err := client.UpdateSource(sourceID, volley.UpdateSourceRequest{
    Name: "Updated Source Name",
    EPS:  &[]int{20}[0],
})
if err != nil {
    log.Fatal(err)
}
```

### Destinations

```go
// List destinations
destinations, err := client.ListDestinations(projectID)
if err != nil {
    log.Fatal(err)
}

// Create destination
dest, err := client.CreateDestination(projectID, volley.CreateDestinationRequest{
    Name: "Production Endpoint",
    URL:  "https://api.example.com/webhooks",
    EPS:  5,
})
if err != nil {
    log.Fatal(err)
}
```

### Connections

```go
// List connections
connections, err := client.GetConnections(projectID)
if err != nil {
    log.Fatal(err)
}

// Create connection
conn, err := client.CreateConnection(projectID, volley.CreateConnectionRequest{
    SourceID:      sourceID,
    DestinationID: destID,
    Status:        "enabled",
    EPS:           5,
    MaxRetries:    3,
})
if err != nil {
    log.Fatal(err)
}
```

### Events

```go
// List events with filters
events, err := client.ListEvents(projectID, &volley.ListEventsOptions{
    Status:   "failed",
    SourceID: &sourceID,
    Limit:    &[]int{50}[0],
    Offset:   &[]int{0}[0],
})
if err != nil {
    log.Fatal(err)
}

// Get event details
event, err := client.GetEvent(requestID)
if err != nil {
    log.Fatal(err)
}

// Replay failed event
result, err := client.ReplayEvent(volley.ReplayEventRequest{
    EventID: "evt_abc123def456",
})
if err != nil {
    log.Fatal(err)
}
```

### Delivery Attempts

```go
// List delivery attempts
attempts, err := client.ListDeliveryAttempts(projectID, &volley.ListDeliveryAttemptsOptions{
    EventID: "evt_abc123",
    Status:  "failed",
    Limit:   &[]int{50}[0],
})
if err != nil {
    log.Fatal(err)
}
```

### Sending Webhooks

```go
// Send a webhook to a source
eventID, err := client.SendWebhook("source_ingestion_id", map[string]interface{}{
    "event": "user.created",
    "data": map[string]interface{}{
        "user_id": "123",
        "email":   "user@example.com",
    },
})
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The SDK returns errors that implement the `error` interface. API errors are returned as `*volley.APIError`:

```go
event, err := client.GetEvent(requestID)
if err != nil {
    if apiErr, ok := err.(*volley.APIError); ok {
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.ErrorMsg, apiErr.Status)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Common HTTP Status Codes

- `200` - Success
- `201` - Created
- `202` - Accepted (webhook queued)
- `400` - Bad Request (validation error)
- `401` - Unauthorized (invalid or missing API token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Rate Limit Exceeded
- `500` - Internal Server Error

For more details on error responses, see the [API Reference - Response Codes](https://docs.volleyhooks.com/api#response-codes).

## Client Options

You can configure the client with various options:

```go
// Custom base URL (for testing)
client := volley.NewClient("token", 
    volley.WithBaseURL("https://api-staging.volleyhooks.com"),
)

// Custom HTTP client
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}
client := volley.NewClient("token",
    volley.WithHTTPClient(httpClient),
)
```

## Additional Resources

### Documentation

- **[Getting Started](https://docs.volleyhooks.com/getting-started)** - Set up your account and create your first webhook
- **[How It Works](https://docs.volleyhooks.com/how-it-works)** - Understand Volley's architecture
- **[Webhooks Guide](https://docs.volleyhooks.com/webhooks)** - Learn about webhook best practices and signature verification
- **[Rate Limiting](https://docs.volleyhooks.com/rate-limiting)** - Configure rate limits for your webhooks
- **[Monitoring](https://docs.volleyhooks.com/monitoring)** - Monitor webhook delivery and performance
- **[Best Practices](https://docs.volleyhooks.com/best-practices)** - Webhook development best practices
- **[FAQ](https://docs.volleyhooks.com/faq)** - Frequently asked questions

### Use Cases

- [Stripe Webhook Localhost Testing](https://docs.volleyhooks.com/use-cases/stripe-webhook-localhost)
- [Retrying Failed Webhooks](https://docs.volleyhooks.com/use-cases/retrying-failed-webhooks)
- [Webhook Event Replay](https://docs.volleyhooks.com/use-cases/webhook-event-replay)
- [Webhook Fan-out](https://docs.volleyhooks.com/use-cases/webhook-fan-out)
- [Multi-Tenant Webhooks](https://docs.volleyhooks.com/use-cases/multi-tenant-webhooks)

## Support

- **Documentation**: [https://docs.volleyhooks.com](https://docs.volleyhooks.com)
- **Console**: [https://app.volleyhooks.com](https://app.volleyhooks.com)
- **Website**: [https://volleyhooks.com](https://volleyhooks.com)

## Testing

The SDK includes comprehensive unit tests and integration tests.

### Running Unit Tests

Unit tests use a mock HTTP server and don't require API credentials:

```bash
go test ./...
```

To run tests with verbose output:

```bash
go test -v ./...
```

### Running Integration Tests

Integration tests make real API calls to the Volley API. You'll need to set your API token:

```bash
export VOLLEY_API_TOKEN="your-api-token"
go test -v -run TestIntegration ./...
```

**Note**: Integration tests are skipped if `VOLLEY_API_TOKEN` is not set or if running in short mode (`go test -short`).

### Test Coverage

To see test coverage:

```bash
go test -cover ./...
```

For detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

When contributing:
1. Add tests for new functionality
2. Ensure all tests pass (`go test ./...`)
3. Update documentation as needed

## License

MIT License - see [LICENSE](LICENSE) file for details.

