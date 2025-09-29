package database

import (
	"database/sql"
	"fmt"
)

// GetArtists retrieves all artists from the database
func (c *Connection) GetArtists() ([]PersistentArtist, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT a.id, a.name, a.name_lower, a.created_at, a.updated_at,
		       COUNT(t.id) as tracks_count
		FROM artists a
		LEFT JOIN tracks t ON a.id = t.artist_id
		GROUP BY a.id, a.name, a.name_lower, a.created_at, a.updated_at
		ORDER BY a.name
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query artists: %w", err)
	}
	defer rows.Close()

	var artists []PersistentArtist
	for rows.Next() {
		var artist PersistentArtist
		err := rows.Scan(
			&artist.ID, &artist.Name, &artist.NameLower,
			&artist.CreatedAt, &artist.UpdatedAt, &artist.TracksCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artist: %w", err)
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

// GetArtistByID retrieves an artist by its ID
func (c *Connection) GetArtistByID(id int64) (*PersistentArtist, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT a.id, a.name, a.name_lower, a.created_at, a.updated_at,
		       COUNT(t.id) as tracks_count
		FROM artists a
		LEFT JOIN tracks t ON a.id = t.artist_id
		WHERE a.id = ?
		GROUP BY a.id, a.name, a.name_lower, a.created_at, a.updated_at
	`

	var artist PersistentArtist
	err := c.db.QueryRow(query, id).Scan(
		&artist.ID, &artist.Name, &artist.NameLower,
		&artist.CreatedAt, &artist.UpdatedAt, &artist.TracksCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artist with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	return &artist, nil
}

// GetTracksByArtistID retrieves all tracks for a specific artist
func (c *Connection) GetTracksByArtistID(artistID int64) ([]PersistentTrack, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT id, file_path, file_name, title, album_name, album_artist_name, 
		       album_id, artist_name, artist_id, image_path, track_number, 
		       txt_lyrics, lrc_lyrics, duration, instrumental, title_lower,
		       created_at, updated_at
		FROM tracks
		WHERE artist_id = ?
		ORDER BY album_name, track_number, title
	`

	rows, err := c.db.Query(query, artistID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks by artist: %w", err)
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
