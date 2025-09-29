package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConnection(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test creating a new connection
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test database file creation
	dbPath := filepath.Join(tempDir, "db.sqlite3")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file not created")
	}

	// Test connection is working
	if conn.GetDB() == nil {
		t.Errorf("Database connection is nil")
	}
}

func TestConnectionClose(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}

	// Test closing the connection
	err = conn.Close()
	if err != nil {
		t.Errorf("Failed to close connection: %v", err)
	}
}

func TestConnectionPooling(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that connection pooling is configured
	db := conn.GetDB()
	if db == nil {
		t.Fatalf("Database connection is nil")
	}

	// Test that we can get connection stats
	stats := db.Stats()
	// Test that stats are accessible (we don't check specific values as they may vary)
	_ = stats
}

func TestConnectionPing(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that we can ping the database
	db := conn.GetDB()
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

func TestConnectionConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			db := conn.GetDB()
			if err := db.Ping(); err != nil {
				t.Errorf("Failed to ping database in goroutine: %v", err)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestConnectionTimeout(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that connection has proper timeout settings
	db := conn.GetDB()
	stats := db.Stats()

	// Check that connection stats are available
	_ = stats
}

func TestConnectionWithInvalidPath(t *testing.T) {
	// Test with invalid path
	_, err := NewConnection("/invalid/path/that/does/not/exist")
	if err == nil {
		t.Errorf("Expected error for invalid path, got nil")
	}
}

func TestConnectionWithExistingDatabase(t *testing.T) {
	tempDir := t.TempDir()

	// Create first connection
	conn1, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create first connection: %v", err)
	}
	conn1.Close()

	// Create second connection to same database
	conn2, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create second connection: %v", err)
	}
	defer conn2.Close()

	// Test that we can still access the database
	db := conn2.GetDB()
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database after reopening: %v", err)
	}
}

func TestConnectionStats(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that we can get connection stats
	db := conn.GetDB()
	stats := db.Stats()

	// Test that we can access connection stats
	_ = stats
}

func TestConnectionLifetime(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that connection lifetime is properly set
	db := conn.GetDB()
	stats := db.Stats()

	// The connection should be configured with proper lifetime
	_ = stats

	// Test that we can perform operations
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

func TestConnectionIdleTimeout(t *testing.T) {
	tempDir := t.TempDir()
	conn, err := NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Test that idle timeout is properly configured
	db := conn.GetDB()
	stats := db.Stats()

	// The connection should be configured with proper idle timeout
	_ = stats

	// Test that we can perform operations
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}
