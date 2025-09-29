package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// GetTracks retrieves all tracks from the database
func (c *Connection) GetTracks() ([]PersistentTrack, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT id, file_path, file_name, title, album_name, album_artist_name, 
		       album_id, artist_name, artist_id, image_path, track_number, 
		       txt_lyrics, lrc_lyrics, duration, instrumental, title_lower,
		       created_at, updated_at
		FROM tracks
		ORDER BY artist_name, album_name, track_number
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}
	defer rows.Close()

	var tracks []PersistentTrack
	for rows.Next() {
		var track PersistentTrack
		err := rows.Scan(
			&track.ID, &track.FilePath, &track.FileName, &track.Title,
			&track.AlbumName, &track.AlbumArtistName, &track.AlbumID,
			&track.ArtistName, &track.ArtistID, &track.ImagePath,
			&track.TrackNumber, &track.TxtLyrics, &track.LrcLyrics,
			&track.Duration, &track.Instrumental, &track.TitleLower,
			&track.CreatedAt, &track.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// GetTrackByID retrieves a track by its ID
func (c *Connection) GetTrackByID(id int64) (*PersistentTrack, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT id, file_path, file_name, title, album_name, album_artist_name, 
		       album_id, artist_name, artist_id, image_path, track_number, 
		       txt_lyrics, lrc_lyrics, duration, instrumental, title_lower,
		       created_at, updated_at
		FROM tracks
		WHERE id = ?
	`

	var track PersistentTrack
	err := c.db.QueryRow(query, id).Scan(
		&track.ID, &track.FilePath, &track.FileName, &track.Title,
		&track.AlbumName, &track.AlbumArtistName, &track.AlbumID,
		&track.ArtistName, &track.ArtistID, &track.ImagePath,
		&track.TrackNumber, &track.TxtLyrics, &track.LrcLyrics,
		&track.Duration, &track.Instrumental, &track.TitleLower,
		&track.CreatedAt, &track.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("track with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

// AddTrack adds a new track to the database
func (c *Connection) AddTrack(track *PersistentTrack) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// First, ensure artist exists
	artistID, err := c.getOrCreateArtist(track.ArtistName)
	if err != nil {
		return fmt.Errorf("failed to get or create artist: %w", err)
	}

	// Then, ensure album exists
	albumID, err := c.getOrCreateAlbum(track.AlbumName, track.AlbumArtistName, track.ImagePath, artistID)
	if err != nil {
		return fmt.Errorf("failed to get or create album: %w", err)
	}

	// Prepare title_lower
	titleLower := strings.ToLower(track.Title)

	query := `
		INSERT INTO tracks (file_path, file_name, title, album_name, album_artist_name,
		                   album_id, artist_name, artist_id, image_path, track_number,
		                   txt_lyrics, lrc_lyrics, duration, instrumental, title_lower,
		                   created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := c.db.Exec(query,
		track.FilePath, track.FileName, track.Title, track.AlbumName, track.AlbumArtistName,
		albumID, track.ArtistName, artistID, track.ImagePath, track.TrackNumber,
		track.TxtLyrics, track.LrcLyrics, track.Duration, track.Instrumental, titleLower,
		now, now,
	)

	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	trackID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get track ID: %w", err)
	}

	track.ID = trackID
	track.AlbumID = albumID
	track.ArtistID = artistID
	track.CreatedAt = now
	track.UpdatedAt = now

	return nil
}

// UpdateTrackSyncedLyrics updates a track's synced lyrics
func (c *Connection) UpdateTrackSyncedLyrics(trackID int64, syncedLyrics, plainLyrics string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `
		UPDATE tracks 
		SET lrc_lyrics = ?, txt_lyrics = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := c.db.Exec(query, syncedLyrics, plainLyrics, time.Now(), trackID)
	if err != nil {
		return fmt.Errorf("failed to update synced lyrics: %w", err)
	}

	return nil
}

// UpdateTrackPlainLyrics updates a track's plain lyrics
func (c *Connection) UpdateTrackPlainLyrics(trackID int64, plainLyrics string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `
		UPDATE tracks 
		SET txt_lyrics = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := c.db.Exec(query, plainLyrics, time.Now(), trackID)
	if err != nil {
		return fmt.Errorf("failed to update plain lyrics: %w", err)
	}

	return nil
}

// UpdateTrackInstrumental marks a track as instrumental
func (c *Connection) UpdateTrackInstrumental(trackID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `
		UPDATE tracks 
		SET instrumental = TRUE, updated_at = ?
		WHERE id = ?
	`

	_, err := c.db.Exec(query, time.Now(), trackID)
	if err != nil {
		return fmt.Errorf("failed to update instrumental status: %w", err)
	}

	return nil
}

// getOrCreateArtist gets an existing artist or creates a new one
func (c *Connection) getOrCreateArtist(artistName string) (int64, error) {
	// Try to find existing artist
	var artistID int64
	err := c.db.QueryRow("SELECT id FROM artists WHERE name = ?", artistName).Scan(&artistID)
	if err == nil {
		return artistID, nil
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query artist: %w", err)
	}

	// Create new artist
	nameLower := strings.ToLower(artistName)
	now := time.Now()
	
	result, err := c.db.Exec(
		"INSERT INTO artists (name, name_lower, created_at, updated_at) VALUES (?, ?, ?, ?)",
		artistName, nameLower, now, now,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create artist: %w", err)
	}

	artistID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get artist ID: %w", err)
	}

	return artistID, nil
}

// getOrCreateAlbum gets an existing album or creates a new one
func (c *Connection) getOrCreateAlbum(albumName string, albumArtistName *string, imagePath *string, artistID int64) (int64, error) {
	// Try to find existing album
	var albumID int64
	query := "SELECT id FROM albums WHERE name = ? AND artist_id = ?"
	err := c.db.QueryRow(query, albumName, artistID).Scan(&albumID)
	if err == nil {
		return albumID, nil
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query album: %w", err)
	}

	// Create new album
	nameLower := strings.ToLower(albumName)
	var albumArtistNameLower *string
	if albumArtistName != nil {
		lower := strings.ToLower(*albumArtistName)
		albumArtistNameLower = &lower
	}

	now := time.Now()
	
	result, err := c.db.Exec(
		"INSERT INTO albums (name, name_lower, artist_name, album_artist_name, album_artist_name_lower, image_path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		albumName, nameLower, "", albumArtistName, albumArtistNameLower, imagePath, now, now,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create album: %w", err)
	}

	albumID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get album ID: %w", err)
	}

	return albumID, nil
}
