// Package sdk provides a minimal Go client for the feature flag service.
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	appID      string
	httpClient *http.Client
}

const DefaultBaseURL = "http://localhost:8080"

func New(baseURL, appID string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		appID:   appID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type publicFlagResponse struct {
	Enabled bool `json:"enabled"`
}

func (c *Client) IsEnabled(ctx context.Context, flagName, userID string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/public/flag/%s/%s", c.baseURL, c.appID, flagName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var out publicFlagResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, err
	}
	return out.Enabled, nil
}
