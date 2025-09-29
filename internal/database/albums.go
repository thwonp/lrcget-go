package database

import (
	"database/sql"
	"fmt"
)

// GetAlbums retrieves all albums from the database
func (c *Connection) GetAlbums() ([]PersistentAlbum, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT a.id, a.name, a.image_path, a.artist_name, a.album_artist_name,
		       a.name_lower, a.album_artist_name_lower, a.created_at, a.updated_at,
		       COUNT(t.id) as tracks_count
		FROM albums a
		LEFT JOIN tracks t ON a.id = t.album_id
		GROUP BY a.id, a.name, a.image_path, a.artist_name, a.album_artist_name,
		         a.name_lower, a.album_artist_name_lower, a.created_at, a.updated_at
		ORDER BY a.artist_name, a.name
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}
	defer rows.Close()

	var albums []PersistentAlbum
	for rows.Next() {
		var album PersistentAlbum
		err := rows.Scan(
			&album.ID, &album.Name, &album.ImagePath, &album.ArtistName,
			&album.AlbumArtistName, &album.NameLower, &album.AlbumArtistNameLower,
			&album.CreatedAt, &album.UpdatedAt, &album.TracksCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}

	return albums, nil
}

// GetAlbumByID retrieves an album by its ID
func (c *Connection) GetAlbumByID(id int64) (*PersistentAlbum, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT a.id, a.name, a.image_path, a.artist_name, a.album_artist_name,
		       a.name_lower, a.album_artist_name_lower, a.created_at, a.updated_at,
		       COUNT(t.id) as tracks_count
		FROM albums a
		LEFT JOIN tracks t ON a.id = t.album_id
		WHERE a.id = ?
		GROUP BY a.id, a.name, a.image_path, a.artist_name, a.album_artist_name,
		         a.name_lower, a.album_artist_name_lower, a.created_at, a.updated_at
	`

	var album PersistentAlbum
	err := c.db.QueryRow(query, id).Scan(
		&album.ID, &album.Name, &album.ImagePath, &album.ArtistName,
		&album.AlbumArtistName, &album.NameLower, &album.AlbumArtistNameLower,
		&album.CreatedAt, &album.UpdatedAt, &album.TracksCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("album with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	return &album, nil
}

// GetTracksByAlbumID retrieves all tracks for a specific album
func (c *Connection) GetTracksByAlbumID(albumID int64) ([]PersistentTrack, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT id, file_path, file_name, title, album_name, album_artist_name, 
		       album_id, artist_name, artist_id, image_path, track_number, 
		       txt_lyrics, lrc_lyrics, duration, instrumental, title_lower,
		       created_at, updated_at
		FROM tracks
		WHERE album_id = ?
		ORDER BY track_number, title
	`

	rows, err := c.db.Query(query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks by album: %w", err)
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
