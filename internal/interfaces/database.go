package interfaces

import (
	"context"
	"lrcget-go/internal/database"
)

// DatabaseInterface defines the interface for database operations
type DatabaseInterface interface {
	// Connection management
	Close() error
	GetDB() interface{} // Returns the underlying database connection
	
	// Track operations
	GetTracks() ([]database.PersistentTrack, error)
	GetTrackByID(id int64) (*database.PersistentTrack, error)
	AddTrack(track *database.PersistentTrack) error
	UpdateTrack(track *database.PersistentTrack) error
	DeleteTrack(id int64) error
	GetTracksByAlbumID(albumID int64) ([]database.PersistentTrack, error)
	GetTracksByArtistID(artistID int64) ([]database.PersistentTrack, error)
	UpdateTrackSyncedLyrics(trackID int64, synced, plain string) error
	UpdateTrackPlainLyrics(trackID int64, plain string) error
	UpdateTrackInstrumental(trackID int64) error
	
	// Album operations
	GetAlbums() ([]database.PersistentAlbum, error)
	GetAlbumByID(id int64) (*database.PersistentAlbum, error)
	AddAlbum(album *database.PersistentAlbum) error
	UpdateAlbum(album *database.PersistentAlbum) error
	DeleteAlbum(id int64) error
	
	// Artist operations
	GetArtists() ([]database.PersistentArtist, error)
	GetArtistByID(id int64) (*database.PersistentArtist, error)
	AddArtist(artist *database.PersistentArtist) error
	UpdateArtist(artist *database.PersistentArtist) error
	DeleteArtist(id int64) error
	
	// Configuration operations
	GetConfig() (*database.PersistentConfig, error)
	UpdateConfig(config *database.PersistentConfig) error
	
	// Directory operations
	GetDirectories() ([]string, error)
	SetDirectories(directories []string) error
	
	// Initialization
	GetInit() (bool, error)
	SetInit(init bool) error
	
	// Migration
	Migrate() error
}

// TransactionInterface defines the interface for database transactions
type TransactionInterface interface {
	Begin() (TransactionInterface, error)
	Commit() error
	Rollback() error
	GetTracks() ([]database.PersistentTrack, error)
	AddTrack(track *database.PersistentTrack) error
	UpdateTrack(track *database.PersistentTrack) error
	DeleteTrack(id int64) error
}

// QueryInterface defines the interface for database queries
type QueryInterface interface {
	Execute(query string, args ...interface{}) (interface{}, error)
	ExecuteWithContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Query(query string, args ...interface{}) (interface{}, error)
	QueryWithContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}
