package lrclib

import "fmt"

// RawResponse represents the raw response from LRCLIB API
type RawResponse struct {
	PlainLyrics  *string  `json:"plainLyrics"`
	SyncedLyrics *string  `json:"syncedLyrics"`
	Instrumental bool     `json:"instrumental"`
	Lang         *string  `json:"lang"`
	Isrc         *string  `json:"isrc"`
	SpotifyID    *string  `json:"spotifyId"`
	Name         *string  `json:"name"`
	AlbumName    *string  `json:"albumName"`
	ArtistName   *string  `json:"artistName"`
	ReleaseDate  *string  `json:"releaseDate"`
	Duration     *float64 `json:"duration"`
}

// Response represents a generic LRCLIB response
type Response interface {
	Type() string
}

// SyncedLyrics represents synced lyrics response
type SyncedLyrics struct {
	Synced string `json:"synced"`
	Plain  string `json:"plain"`
}

func (s SyncedLyrics) Type() string { return "synced" }

// UnsyncedLyrics represents unsynced lyrics response
type UnsyncedLyrics struct {
	Plain string `json:"plain"`
}

func (u UnsyncedLyrics) Type() string { return "unsynced" }

// Instrumental represents instrumental response
type Instrumental struct{}

func (i Instrumental) Type() string { return "instrumental" }

// None represents no lyrics found
type None struct{}

func (n None) Type() string { return "none" }

// ChallengeResponse represents a challenge request response
type ChallengeResponse struct {
	Prefix string `json:"prefix"`
	Target string `json:"target"`
}

// SearchResult represents a search result
type SearchResult struct {
	ID           int64   `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	SyncedLyrics *string `json:"syncedLyrics"`
	PlainLyrics  *string `json:"plainLyrics"`
	Instrumental bool    `json:"instrumental"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Data []SearchResult `json:"data"`
}

// PublishRequest represents a publish request
type PublishRequest struct {
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	SyncedLyrics *string `json:"syncedLyrics"`
	PlainLyrics  *string `json:"plainLyrics"`
	Instrumental bool    `json:"instrumental"`
}

// PublishResponse represents a publish response
type PublishResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

// FlagRequest represents a flag request
type FlagRequest struct {
	TrackID int64  `json:"trackId"`
	Reason  string `json:"reason"`
}

// APIError represents an API error
type APIError struct {
	StatusCode *int   `json:"statusCode"`
	ErrorType  string `json:"error"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
}
