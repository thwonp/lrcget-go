package lrclib

import (
	"context"
	"fmt"
	"net/http"
)

// FlagLyrics flags lyrics on the LRCLIB API
func (c *Client) FlagLyrics(ctx context.Context, req FlagRequest) error {
	apiURL := fmt.Sprintf("%s/api/flag", c.baseURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("Content-Type", "application/json")
	
	if err := encodeJSONRequest(httpReq, req); err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}
	
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	return nil
}
