package database

import (
	"time"
)

// PersistentTrack represents a track in the database
type PersistentTrack struct {
	ID                 int64   `json:"id" db:"id"`
	FilePath           string  `json:"file_path" db:"file_path"`
	FileName           string  `json:"file_name" db:"file_name"`
	Title              string  `json:"title" db:"title"`
	AlbumName          string  `json:"album_name" db:"album_name"`
	AlbumArtistName    *string `json:"album_artist_name" db:"album_artist_name"`
	AlbumID            int64   `json:"album_id" db:"album_id"`
	ArtistName         string  `json:"artist_name" db:"artist_name"`
	ArtistID           int64   `json:"artist_id" db:"artist_id"`
	ImagePath          *string `json:"image_path" db:"image_path"`
	TrackNumber        *int64  `json:"track_number" db:"track_number"`
	TxtLyrics          *string `json:"txt_lyrics" db:"txt_lyrics"`
	LrcLyrics          *string `json:"lrc_lyrics" db:"lrc_lyrics"`
	Duration           float64 `json:"duration" db:"duration"`
	Instrumental       bool    `json:"instrumental" db:"instrumental"`
	TitleLower         *string `json:"title_lower" db:"title_lower"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// PersistentAlbum represents an album in the database
type PersistentAlbum struct {
	ID             int64   `json:"id" db:"id"`
	Name           string  `json:"name" db:"name"`
	ImagePath      *string `json:"image_path" db:"image_path"`
	ArtistName     string  `json:"artist_name" db:"artist_name"`
	AlbumArtistName *string `json:"album_artist_name" db:"album_artist_name"`
	NameLower      *string `json:"name_lower" db:"name_lower"`
	AlbumArtistNameLower *string `json:"album_artist_name_lower" db:"album_artist_name_lower"`
	TracksCount    int64   `json:"tracks_count" db:"tracks_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// PersistentArtist represents an artist in the database
type PersistentArtist struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	NameLower   *string `json:"name_lower" db:"name_lower"`
	TracksCount int64  `json:"tracks_count" db:"tracks_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PersistentConfig represents application configuration
type PersistentConfig struct {
	ID                           int64  `json:"id" db:"id"`
	SkipTracksWithSyncedLyrics   bool   `json:"skip_tracks_with_synced_lyrics" db:"skip_tracks_with_synced_lyrics"`
	SkipTracksWithPlainLyrics    bool   `json:"skip_tracks_with_plain_lyrics" db:"skip_tracks_with_plain_lyrics"`
	ShowLineCount                bool   `json:"show_line_count" db:"show_line_count"`
	TryEmbedLyrics               bool   `json:"try_embed_lyrics" db:"try_embed_lyrics"`
	ThemeMode                    string `json:"theme_mode" db:"theme_mode"`
	LrclibInstance               string `json:"lrclib_instance" db:"lrclib_instance"`
	CreatedAt                    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" db:"updated_at"`
}

// PersistentDirectory represents a directory in the database
type PersistentDirectory struct {
	ID        int64     `json:"id" db:"id"`
	Path      string    `json:"path" db:"path"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PersistentInit represents library initialization status
type PersistentInit struct {
	ID        int64     `json:"id" db:"id"`
	Initialized bool    `json:"initialized" db:"initialized"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
