package lrclib

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SearchLyrics searches for lyrics using the LRCLIB API
func (c *Client) SearchLyrics(ctx context.Context, title, artist, album, query string) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	if title != "" {
		params.Set("track_name", title)
	}
	if artist != "" {
		params.Set("artist_name", artist)
	}
	if album != "" {
		params.Set("album_name", album)
	}
	
	apiURL := fmt.Sprintf("%s/api/search?%s", c.baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	var searchResp SearchResponse
	if err := decodeJSONResponse(resp, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &searchResp, nil
}