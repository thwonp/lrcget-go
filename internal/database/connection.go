package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const CurrentDBVersion = 7

// Connection represents a database connection
type Connection struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewConnection creates a new database connection with connection pooling
func NewConnection(dataDir string) (*Connection, error) {
	// Create data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database
	dbPath := filepath.Join(dataDir, "db.sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pooling
	db.SetMaxOpenConns(25)        // Maximum number of open connections
	db.SetMaxIdleConns(5)         // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute)  // Maximum connection lifetime
	db.SetConnMaxIdleTime(1 * time.Minute)   // Maximum idle time

	conn := &Connection{db: db}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := conn.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return conn, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.db.Close()
}

// GetDB returns the underlying database connection (for internal use)
func (c *Connection) GetDB() *sql.DB {
	return c.db
}
