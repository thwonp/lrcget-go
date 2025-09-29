package lrclib

import (
	"context"
	"fmt"
	"net/http"
)

// PublishLyrics publishes lyrics to the LRCLIB API
func (c *Client) PublishLyrics(ctx context.Context, req PublishRequest) (*PublishResponse, error) {
	apiURL := fmt.Sprintf("%s/api/publish", c.baseURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("Content-Type", "application/json")
	
	if err := encodeJSONRequest(httpReq, req); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	var publishResp PublishResponse
	if err := decodeJSONResponse(resp, &publishResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &publishResp, nil
}