# Testing Guide

This document explains how to run and write tests for the Volley Go SDK.

## Test Structure

The SDK includes two types of tests:

1. **Unit Tests** - Use mock HTTP servers, no API credentials required
2. **Integration Tests** - Make real API calls, require API token

## Running Tests

### Unit Tests (No API Token Required)

Run all unit tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

Run specific test:

```bash
go test -v -run TestListOrganizations
```

### Integration Tests (Requires API Token)

Integration tests make real API calls to test the SDK against the actual Volley API.

**Setup:**

1. Get your API token from [Volley Console](https://app.volleyhooks.com) → Settings → Account → API Token
2. Set the environment variable:

```bash
export VOLLEY_API_TOKEN="your-api-token"
```

**Run integration tests:**

```bash
go test -v -run TestIntegration ./...
```

**Note**: Integration tests are automatically skipped if:
- `VOLLEY_API_TOKEN` is not set
- Running in short mode (`go test -short`)

## Test Coverage

View test coverage:

```bash
go test -cover ./...
```

Generate detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This will open a browser showing which lines are covered by tests.

## Test Files

- `client_test.go` - Client initialization and configuration tests
- `organizations_test.go` - Organization API tests
- `projects_test.go` - Project API tests
- `sources_test.go` - Source API tests
- `events_test.go` - Event and replay API tests
- `integration_test.go` - Real API integration tests

## Writing New Tests

### Unit Test Example

```go
func TestMyNewFeature(t *testing.T) {
    // Create a test server
    server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
        // Verify request
        if r.Method != http.MethodGet {
            t.Errorf("Expected GET, got %s", r.Method)
        }
        
        // Return mock response
        response := map[string]interface{}{
            "data": "test",
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })
    defer server.Close()
    
    // Create client with test server URL
    client := volley.NewClient("test-token", volley.WithBaseURL(server.URL))
    
    // Test your feature
    result, err := client.MyNewFeature()
    if err != nil {
        t.Fatalf("MyNewFeature failed: %v", err)
    }
    
    // Assertions
    if result != "expected" {
        t.Errorf("Expected 'expected', got %s", result)
    }
}
```

### Integration Test Example

```go
func TestIntegrationMyFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    apiToken := os.Getenv("VOLLEY_API_TOKEN")
    if apiToken == "" {
        t.Skip("Skipping integration test: VOLLEY_API_TOKEN not set")
    }
    
    client := volley.NewClient(apiToken)
    
    // Test against real API
    result, err := client.MyNewFeature()
    if err != nil {
        t.Fatalf("MyNewFeature failed: %v", err)
    }
    
    t.Logf("Result: %v", result)
}
```

## Best Practices

1. **Always test error cases** - Test both success and failure scenarios
2. **Verify request details** - Check HTTP method, path, headers, body
3. **Use descriptive test names** - `TestListOrganizations` is better than `Test1`
4. **Clean up resources** - Use `defer` for cleanup
5. **Skip integration tests** - When API token is not available or in short mode
6. **Use table-driven tests** - For testing multiple scenarios

## Continuous Integration

For CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run unit tests
  run: go test ./...

- name: Run integration tests
  env:
    VOLLEY_API_TOKEN: ${{ secrets.VOLLEY_API_TOKEN }}
  run: go test -v -run TestIntegration ./...
```

## Troubleshooting

### Tests fail with "connection refused"
- Make sure test server is created before making requests
- Check that `defer server.Close()` is called

### Integration tests skipped
- Verify `VOLLEY_API_TOKEN` environment variable is set
- Check that you're not running with `-short` flag

### Test coverage is low
- Add tests for error cases
- Test edge cases and boundary conditions
- Test all public methods

