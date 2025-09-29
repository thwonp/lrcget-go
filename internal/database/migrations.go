package database

import (
	"fmt"
)

// Migrate runs database migrations
func (c *Connection) Migrate() error {
	// Get current version
	var version int
	err := c.db.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to get database version: %w", err)
	}

	if version < CurrentDBVersion {
		return c.upgradeDatabase(version)
	}

	return nil
}

// upgradeDatabase upgrades the database from the given version to the current version
func (c *Connection) upgradeDatabase(fromVersion int) error {
	fmt.Printf("Existing database version: %d\n", fromVersion)

	if fromVersion <= 0 {
		fmt.Println("Migrate database version 1...")
		return c.migrateToVersion1()
	}

	if fromVersion <= 1 {
		fmt.Println("Migrate database version 2...")
		return c.migrateToVersion2()
	}

	if fromVersion <= 2 {
		fmt.Println("Migrate database version 3...")
		return c.migrateToVersion3()
	}

	if fromVersion <= 3 {
		fmt.Println("Migrate database version 4...")
		return c.migrateToVersion4()
	}

	if fromVersion <= 4 {
		fmt.Println("Migrate database version 5...")
		return c.migrateToVersion5()
	}

	if fromVersion <= 5 {
		fmt.Println("Migrate database version 6...")
		return c.migrateToVersion6()
	}

	if fromVersion <= 6 {
		fmt.Println("Migrate database version 7...")
		return c.migrateToVersion7()
	}

	return nil
}

// migrateToVersion1 creates the initial database schema
func (c *Connection) migrateToVersion1() error {
	// Set journal mode (must be done outside transaction)
	_, err := c.db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return fmt.Errorf("failed to set journal mode: %w", err)
	}

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Set user version
	_, err = tx.Exec("PRAGMA user_version = 1")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	// Create tables
	schema := `
	CREATE TABLE directories (
		id INTEGER PRIMARY KEY,
		path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE library_data (
		id INTEGER PRIMARY KEY,
		init BOOLEAN,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE config_data (
		id INTEGER PRIMARY KEY,
		skip_not_needed_tracks BOOLEAN,
		try_embed_lyrics BOOLEAN,
		skip_tracks_with_synced_lyrics BOOLEAN DEFAULT 0,
		skip_tracks_with_plain_lyrics BOOLEAN DEFAULT 0,
		show_line_count BOOLEAN DEFAULT 1,
		theme_mode TEXT DEFAULT 'system',
		lrclib_instance TEXT DEFAULT 'https://lrclib.net',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE artists (
		id INTEGER PRIMARY KEY,
		name TEXT,
		name_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE albums (
		id INTEGER PRIMARY KEY,
		name TEXT,
		artist_id INTEGER,
		image_path TEXT,
		artist_name TEXT,
		album_artist_name TEXT,
		name_lower TEXT,
		album_artist_name_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(artist_id) REFERENCES artists(id)
	);

	CREATE TABLE tracks (
		id INTEGER PRIMARY KEY,
		file_path TEXT,
		file_name TEXT,
		title TEXT,
		album_name TEXT,
		artist_name TEXT,
		album_artist_name TEXT,
		album_id INTEGER,
		artist_id INTEGER,
		image_path TEXT,
		track_number INTEGER,
		txt_lyrics TEXT,
		duration FLOAT,
		lrc_lyrics TEXT,
		instrumental BOOLEAN,
		title_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(artist_id) REFERENCES artists(id),
		FOREIGN KEY(album_id) REFERENCES albums(id)
	);

	INSERT INTO library_data (init) VALUES (0);
	INSERT INTO config_data (skip_not_needed_tracks, try_embed_lyrics, skip_tracks_with_synced_lyrics, skip_tracks_with_plain_lyrics, show_line_count, theme_mode, lrclib_instance) VALUES (1, 0, 1, 0, 1, 'system', 'https://lrclib.net');
	`

	_, err = tx.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return tx.Commit()
}

// migrateToVersion2 adds txt_lyrics column and indexes
func (c *Connection) migrateToVersion2() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return fmt.Errorf("failed to set journal mode: %w", err)
	}

	_, err = tx.Exec("PRAGMA user_version = 2")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	// txt_lyrics column is already in the initial schema, no need to add it here

	_, err = tx.Exec("CREATE INDEX idx_tracks_title ON tracks(title)")
	if err != nil {
		return fmt.Errorf("failed to create tracks title index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_albums_name ON albums(name)")
	if err != nil {
		return fmt.Errorf("failed to create albums name index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_artists_name ON artists(name)")
	if err != nil {
		return fmt.Errorf("failed to create artists name index: %w", err)
	}

	return tx.Commit()
}

// migrateToVersion3 adds instrumental column (already in initial schema)
func (c *Connection) migrateToVersion3() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA user_version = 3")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	// instrumental column is already in the initial schema, no need to add it here

	return tx.Commit()
}

// migrateToVersion4 adds lowercase columns and indexes
func (c *Connection) migrateToVersion4() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA user_version = 4")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD title_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add title_lower column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD name_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add name_lower column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE artists ADD name_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add name_lower column: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_tracks_title_lower ON tracks(title_lower)")
	if err != nil {
		return fmt.Errorf("failed to create tracks title_lower index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_albums_name_lower ON albums(name_lower)")
	if err != nil {
		return fmt.Errorf("failed to create albums name_lower index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_artists_name_lower ON artists(name_lower)")
	if err != nil {
		return fmt.Errorf("failed to create artists name_lower index: %w", err)
	}

	return tx.Commit()
}

// migrateToVersion5 adds track_number, album_artist_name, and config columns
func (c *Connection) migrateToVersion5() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA user_version = 5")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD track_number INTEGER")
	if err != nil {
		return fmt.Errorf("failed to add track_number column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD album_artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD album_artist_name_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name_lower column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD theme_mode TEXT DEFAULT 'auto'")
	if err != nil {
		return fmt.Errorf("failed to add theme_mode column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD lrclib_instance TEXT DEFAULT 'https://lrclib.net'")
	if err != nil {
		return fmt.Errorf("failed to add lrclib_instance column: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_albums_album_artist_name_lower ON albums(album_artist_name_lower)")
	if err != nil {
		return fmt.Errorf("failed to create albums album_artist_name_lower index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX idx_tracks_track_number ON tracks(track_number)")
	if err != nil {
		return fmt.Errorf("failed to create tracks track_number index: %w", err)
	}

	// Clear existing data and reset initialization
	_, err = tx.Exec("DELETE FROM tracks WHERE 1")
	if err != nil {
		return fmt.Errorf("failed to clear tracks: %w", err)
	}

	_, err = tx.Exec("DELETE FROM albums WHERE 1")
	if err != nil {
		return fmt.Errorf("failed to clear albums: %w", err)
	}

	_, err = tx.Exec("DELETE FROM artists WHERE 1")
	if err != nil {
		return fmt.Errorf("failed to clear artists: %w", err)
	}

	_, err = tx.Exec("UPDATE library_data SET init = 0 WHERE 1")
	if err != nil {
		return fmt.Errorf("failed to reset library initialization: %w", err)
	}

	return tx.Commit()
}

// migrateToVersion6 adds skip columns and renames config columns
func (c *Connection) migrateToVersion6() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA user_version = 6")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD skip_tracks_with_synced_lyrics BOOLEAN DEFAULT 0")
	if err != nil {
		return fmt.Errorf("failed to add skip_tracks_with_synced_lyrics column: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD skip_tracks_with_plain_lyrics BOOLEAN DEFAULT 0")
	if err != nil {
		return fmt.Errorf("failed to add skip_tracks_with_plain_lyrics column: %w", err)
	}

	// Add created_at and updated_at columns to all tables
	_, err = tx.Exec("ALTER TABLE tracks ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to albums: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to albums: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE artists ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to artists: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE artists ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to artists: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE directories ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to directories: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE directories ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to directories: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE library_data ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to library_data: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE library_data ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to library_data: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD created_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add created_at to config_data: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD updated_at DATETIME")
	if err != nil {
		return fmt.Errorf("failed to add updated_at to config_data: %w", err)
	}

	// Add missing columns to tracks table
	_, err = tx.Exec("ALTER TABLE tracks ADD album_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_name to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add artist_name to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD album_artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name to tracks: %w", err)
	}

	// Add missing columns to albums table
	_, err = tx.Exec("ALTER TABLE albums ADD artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add artist_name to albums: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD album_artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name to albums: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE albums ADD album_artist_name_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name_lower to albums: %w", err)
	}

	// Add missing columns to tracks table
	_, err = tx.Exec("ALTER TABLE tracks ADD image_path TEXT")
	if err != nil {
		return fmt.Errorf("failed to add image_path to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD track_number INTEGER")
	if err != nil {
		return fmt.Errorf("failed to add track_number to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD txt_lyrics TEXT")
	if err != nil {
		return fmt.Errorf("failed to add txt_lyrics to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD instrumental BOOLEAN")
	if err != nil {
		return fmt.Errorf("failed to add instrumental to tracks: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE tracks ADD title_lower TEXT")
	if err != nil {
		return fmt.Errorf("failed to add title_lower to tracks: %w", err)
	}

	// Add missing columns that might not exist in older databases
	_, err = tx.Exec("ALTER TABLE tracks ADD album_artist_name TEXT")
	if err != nil {
		return fmt.Errorf("failed to add album_artist_name to tracks: %w", err)
	}

	_, err = tx.Exec("UPDATE config_data SET skip_tracks_with_synced_lyrics = skip_not_needed_tracks")
	if err != nil {
		return fmt.Errorf("failed to migrate skip_not_needed_tracks: %w", err)
	}

	// SQLite doesn't support DROP COLUMN, so we'll leave the old column
	// In a real migration, you'd need to recreate the table

	return tx.Commit()
}

// migrateToVersion7 adds show_line_count column
func (c *Connection) migrateToVersion7() error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("PRAGMA user_version = 7")
	if err != nil {
		return fmt.Errorf("failed to set user version: %w", err)
	}

	_, err = tx.Exec("ALTER TABLE config_data ADD show_line_count BOOLEAN DEFAULT 1")
	if err != nil {
		return fmt.Errorf("failed to add show_line_count column: %w", err)
	}

	return tx.Commit()
}

// createInitialSchema creates the complete current schema (for new installations)
func (c *Connection) createInitialSchema() error {
	schema := `
	PRAGMA journal_mode = WAL;
	
	CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL,
		file_name TEXT NOT NULL,
		title TEXT NOT NULL,
		album_name TEXT NOT NULL,
		album_artist_name TEXT,
		album_id INTEGER NOT NULL,
		artist_name TEXT NOT NULL,
		artist_id INTEGER NOT NULL,
		image_path TEXT,
		track_number INTEGER,
		txt_lyrics TEXT,
		lrc_lyrics TEXT,
		duration REAL NOT NULL,
		instrumental BOOLEAN NOT NULL DEFAULT FALSE,
		title_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (album_id) REFERENCES albums(id),
		FOREIGN KEY (artist_id) REFERENCES artists(id)
	);
	
	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		image_path TEXT,
		artist_name TEXT NOT NULL,
		album_artist_name TEXT,
		name_lower TEXT,
		album_artist_name_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		name_lower TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS directories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS config_data (
		id INTEGER PRIMARY KEY,
		skip_tracks_with_synced_lyrics BOOLEAN NOT NULL DEFAULT FALSE,
		skip_tracks_with_plain_lyrics BOOLEAN NOT NULL DEFAULT FALSE,
		show_line_count BOOLEAN NOT NULL DEFAULT TRUE,
		try_embed_lyrics BOOLEAN NOT NULL DEFAULT FALSE,
		theme_mode TEXT NOT NULL DEFAULT 'system',
		lrclib_instance TEXT NOT NULL DEFAULT 'https://lrclib.net',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS library_data (
		id INTEGER PRIMARY KEY,
		init BOOLEAN NOT NULL DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	-- Create indexes
	CREATE INDEX IF NOT EXISTS idx_tracks_title ON tracks(title);
	CREATE INDEX IF NOT EXISTS idx_tracks_title_lower ON tracks(title_lower);
	CREATE INDEX IF NOT EXISTS idx_tracks_track_number ON tracks(track_number);
	CREATE INDEX IF NOT EXISTS idx_albums_name ON albums(name);
	CREATE INDEX IF NOT EXISTS idx_albums_name_lower ON albums(name_lower);
	CREATE INDEX IF NOT EXISTS idx_albums_album_artist_name_lower ON albums(album_artist_name_lower);
	CREATE INDEX IF NOT EXISTS idx_artists_name ON artists(name);
	CREATE INDEX IF NOT EXISTS idx_artists_name_lower ON artists(name_lower);
	
	-- Insert default data
	INSERT OR IGNORE INTO library_data (id, init) VALUES (1, 0);
	INSERT OR IGNORE INTO config_data (id, skip_tracks_with_synced_lyrics, skip_tracks_with_plain_lyrics, show_line_count, try_embed_lyrics, theme_mode, lrclib_instance) 
	VALUES (1, 1, 0, 1, 0, 'system', 'https://lrclib.net');
	
	PRAGMA user_version = 7;
	`

	_, err := c.db.Exec(schema)
	return err
}
