package volley

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendWebhook sends a webhook to a source
// The sourceID is the ingestion ID provided when you create a source
func (c *Client) SendWebhook(sourceID string, payload interface{}) (string, error) {
	path := fmt.Sprintf("/hook/%s", sourceID)

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	reqURL := c.baseURL + path
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// Note: Authentication is optional for webhook ingestion endpoints
	// If the source has auth configured, include it here

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to get event_id
	var result struct {
		EventID string `json:"event_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.EventID, nil
}

