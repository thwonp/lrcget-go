package interfaces

import (
	"context"
	"lrcget-go/internal/lrclib"
)

// LRCLibInterface defines the interface for LRCLIB API operations
type LRCLibInterface interface {
	// Client configuration
	SetBaseURL(baseURL string)
	GetBaseURL() string

	// Lyrics operations
	GetLyrics(ctx context.Context, title, artist, album string, duration float64) (interface{}, error)
	SearchLyrics(ctx context.Context, title, artist, album, query string) (*lrclib.SearchResponse, error)
	PublishLyrics(ctx context.Context, req lrclib.PublishRequest) (*lrclib.PublishResponse, error)
	FlagLyrics(ctx context.Context, req lrclib.FlagRequest) error

	// Challenge operations
	RequestChallenge(ctx context.Context) (*lrclib.ChallengeResponse, error)
	SolveChallenge(ctx context.Context, challenge string) (interface{}, error)

	// Flag operations
	GetFlag(ctx context.Context, trackID int64) (interface{}, error)
	SetFlag(ctx context.Context, trackID int64, flag string) error

	// Health check
	HealthCheck(ctx context.Context) error
}

// LyricsProviderInterface defines the interface for lyrics providers
type LyricsProviderInterface interface {
	// Search operations
	Search(ctx context.Context, query string) ([]LyricsResult, error)
	SearchByTrack(ctx context.Context, title, artist, album string) ([]LyricsResult, error)
	SearchByArtist(ctx context.Context, artist string) ([]LyricsResult, error)
	SearchByAlbum(ctx context.Context, album string) ([]LyricsResult, error)

	// Get operations
	GetLyrics(ctx context.Context, id string) (*LyricsResult, error)
	GetLyricsByTrack(ctx context.Context, title, artist, album string) (*LyricsResult, error)

	// Publish operations
	PublishLyrics(ctx context.Context, lyrics *LyricsResult) error
	UpdateLyrics(ctx context.Context, id string, lyrics *LyricsResult) error
	DeleteLyrics(ctx context.Context, id string) error

	// Validation
	ValidateLyrics(lyrics string) error
	ValidateTrack(title, artist, album string) error
}

// LyricsResult represents a lyrics search result
type LyricsResult interface {
	GetID() string
	GetTitle() string
	GetArtist() string
	GetAlbum() string
	GetDuration() float64
	GetSyncedLyrics() string
	GetPlainLyrics() string
	GetInstrumental() bool
	GetProvider() string
	GetQuality() int
}

// ChallengeInterface defines the interface for challenge operations
type ChallengeInterface interface {
	Request(ctx context.Context) (*lrclib.ChallengeResponse, error)
	Solve(ctx context.Context, challenge string) (interface{}, error)
	Validate(ctx context.Context, solution string) error
}

// FlagInterface defines the interface for flag operations
type FlagInterface interface {
	Get(ctx context.Context, trackID int64) (interface{}, error)
	Set(ctx context.Context, trackID int64, flag string) error
	Remove(ctx context.Context, trackID int64) error
	List(ctx context.Context) ([]interface{}, error)
}
