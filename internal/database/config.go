package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetConfig retrieves the application configuration
func (c *Connection) GetConfig() (*PersistentConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `
		SELECT id, skip_tracks_with_synced_lyrics, skip_tracks_with_plain_lyrics,
		       show_line_count, try_embed_lyrics, theme_mode, lrclib_instance,
		       created_at, updated_at
		FROM config_data
		WHERE id = 1
	`

	var config PersistentConfig
	err := c.db.QueryRow(query).Scan(
		&config.ID, &config.SkipTracksWithSyncedLyrics, &config.SkipTracksWithPlainLyrics,
		&config.ShowLineCount, &config.TryEmbedLyrics, &config.ThemeMode, &config.LrclibInstance,
		&config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return default config if none exists
			return &PersistentConfig{
				ID:                           1,
				SkipTracksWithSyncedLyrics:   true,
				SkipTracksWithPlainLyrics:    false,
				ShowLineCount:                true,
				TryEmbedLyrics:               false,
				ThemeMode:                    "system",
				LrclibInstance:               "https://lrclib.net",
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return &config, nil
}

// UpdateConfig updates the application configuration
func (c *Connection) UpdateConfig(config *PersistentConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `
		UPDATE config_data 
		SET skip_tracks_with_synced_lyrics = ?, skip_tracks_with_plain_lyrics = ?,
		    show_line_count = ?, try_embed_lyrics = ?, theme_mode = ?, 
		    lrclib_instance = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := c.db.Exec(query,
		config.SkipTracksWithSyncedLyrics, config.SkipTracksWithPlainLyrics,
		config.ShowLineCount, config.TryEmbedLyrics, config.ThemeMode,
		config.LrclibInstance, time.Now(), config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	config.UpdatedAt = time.Now()
	return nil
}

// GetInit retrieves the library initialization status
func (c *Connection) GetInit() (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `SELECT init FROM library_data WHERE id = 1`

	var initialized bool
	err := c.db.QueryRow(query).Scan(&initialized)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to get init status: %w", err)
	}

	return initialized, nil
}

// SetInit sets the library initialization status
func (c *Connection) SetInit(initialized bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `UPDATE library_data SET init = ?, updated_at = ? WHERE id = 1`

	_, err := c.db.Exec(query, initialized, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set init status: %w", err)
	}

	return nil
}

// GetDirectories retrieves all directories from the database
func (c *Connection) GetDirectories() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := `SELECT path FROM directories ORDER BY path`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query directories: %w", err)
	}
	defer rows.Close()

	var directories []string
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory: %w", err)
		}
		directories = append(directories, path)
	}

	return directories, nil
}

// SetDirectories sets the directories in the database
func (c *Connection) SetDirectories(directories []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Start transaction
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing directories
	_, err = tx.Exec("DELETE FROM directories")
	if err != nil {
		return fmt.Errorf("failed to clear directories: %w", err)
	}

	// Insert new directories
	now := time.Now()
	for _, path := range directories {
		_, err = tx.Exec(
			"INSERT INTO directories (path, created_at, updated_at) VALUES (?, ?, ?)",
			path, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert directory: %w", err)
		}
	}

	return tx.Commit()
}
