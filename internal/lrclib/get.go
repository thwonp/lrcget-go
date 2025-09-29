package lrclib

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetLyrics retrieves lyrics for a track
func (c *Client) GetLyrics(ctx context.Context, title, album, artist string, duration float64) (Response, error) {
	params := url.Values{}
	params.Set("track_name", title)
	params.Set("artist_name", artist)
	params.Set("album_name", album)
	params.Set("duration", strconv.FormatFloat(duration, 'f', 2, 64))
	
	apiURL := fmt.Sprintf("%s/api/get?%s", c.baseURL, params.Encode())
	
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
	
	if resp.StatusCode == http.StatusNotFound {
		return None{}, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	var rawResp RawResponse
	if err := decodeJSONResponse(resp, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return c.convertResponse(rawResp), nil
}

// GetLyricsByID retrieves lyrics by track ID
func (c *Client) GetLyricsByID(ctx context.Context, trackID int64) (Response, error) {
	apiURL := fmt.Sprintf("%s/api/get/%d", c.baseURL, trackID)
	
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
	
	if resp.StatusCode == http.StatusNotFound {
		return None{}, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	var rawResp RawResponse
	if err := decodeJSONResponse(resp, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return c.convertResponse(rawResp), nil
}

// convertResponse converts a raw response to the appropriate response type
func (c *Client) convertResponse(raw RawResponse) Response {
	if raw.SyncedLyrics != nil {
		plain := raw.PlainLyrics
		if plain == nil {
			plain = c.stripTimestamp(*raw.SyncedLyrics)
		}
		return SyncedLyrics{
			Synced: *raw.SyncedLyrics,
			Plain:  *plain,
		}
	}
	
	if raw.PlainLyrics != nil {
		return UnsyncedLyrics{Plain: *raw.PlainLyrics}
	}
	
	if raw.Instrumental {
		return Instrumental{}
	}
	
	return None{}
}

// stripTimestamp strips timestamps from synced lyrics to create plain lyrics
func (c *Client) stripTimestamp(syncedLyrics string) *string {
	// This is a simplified implementation
	// In the real version, you'd use regex to remove [mm:ss.xx] patterns
	// For now, just return the synced lyrics as plain
	return &syncedLyrics
}
