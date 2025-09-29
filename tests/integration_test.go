package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"lrcget-go/internal/app"
	"lrcget-go/internal/database"
)

func TestFullWorkflow(t *testing.T) {
	// Clear any existing database
	os.RemoveAll(os.Getenv("HOME") + "/.lrcget")

	// Initialize app
	app := app.NewApp()
	ctx := context.Background()
	app.OnStartup(ctx)
	defer app.OnShutdown(ctx)

	// Create test tracks directly in the database
	createTestTracksDirectly(t, app)

	// Test track retrieval
	tracks, err := app.GetTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 3 {
		t.Errorf("Expected 3 tracks, got %d", len(tracks))
	}

	// Test album operations
	albums, err := app.GetAlbums()
	if err != nil {
		t.Fatalf("Failed to get albums: %v", err)
	}

	if len(albums) == 0 {
		t.Errorf("Expected at least one album")
	}

	// Test artist operations
	artists, err := app.GetArtists()
	if err != nil {
		t.Fatalf("Failed to get artists: %v", err)
	}

	if len(artists) == 0 {
		t.Errorf("Expected at least one artist")
	}

	// Test track by album
	if len(albums) > 0 {
		albumTracks, err := app.GetTracksByAlbum(albums[0].ID)
		if err != nil {
			t.Fatalf("Failed to get tracks by album: %v", err)
		}

		if len(albumTracks) == 0 {
			t.Errorf("Expected at least one track in album")
		}
	}

	// Test track by artist
	if len(artists) > 0 {
		artistTracks, err := app.GetTracksByArtist(artists[0].ID)
		if err != nil {
			t.Fatalf("Failed to get tracks by artist: %v", err)
		}

		if len(artistTracks) == 0 {
			t.Errorf("Expected at least one track by artist")
		}
	}
}

func TestDatabaseOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create database connection
	conn, err := database.NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer conn.Close()

	// Test database initialization
	init, err := conn.GetInit()
	if err != nil {
		t.Fatalf("Failed to get init status: %v", err)
	}

	if init {
		t.Errorf("Expected init to be false initially")
	}

	// Test setting init status
	err = conn.SetInit(true)
	if err != nil {
		t.Fatalf("Failed to set init status: %v", err)
	}

	// Test getting init status again
	init, err = conn.GetInit()
	if err != nil {
		t.Fatalf("Failed to get init status: %v", err)
	}

	if !init {
		t.Errorf("Expected init to be true after setting")
	}

	// Test directory operations
	directories := []string{"/tmp/music1", "/tmp/music2"}
	err = conn.SetDirectories(directories)
	if err != nil {
		t.Fatalf("Failed to set directories: %v", err)
	}

	retrievedDirs, err := conn.GetDirectories()
	if err != nil {
		t.Fatalf("Failed to get directories: %v", err)
	}

	if len(retrievedDirs) != len(directories) {
		t.Errorf("Expected %d directories, got %d", len(directories), len(retrievedDirs))
	}
}

func TestConfigurationOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create database connection
	conn, err := database.NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer conn.Close()

	// Test getting default config
	config, err := conn.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config == nil {
		t.Errorf("Expected config to be non-nil")
	}

	// Test updating config
	config.LrclibInstance = "https://test.lrclib.net"
	err = conn.UpdateConfig(config)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Test getting updated config
	updatedConfig, err := conn.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get updated config: %v", err)
	}

	if updatedConfig.LrclibInstance != "https://test.lrclib.net" {
		t.Errorf("Expected updated config to have new LRCLIB instance")
	}
}

func TestTrackOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create database connection
	conn, err := database.NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer conn.Close()

	// Create test track
	track := &database.PersistentTrack{
		FilePath:   "/tmp/test.mp3",
		FileName:   "test.mp3",
		Title:      "Test Song",
		AlbumName:  "Test Album",
		ArtistName: "Test Artist",
		Duration:   180.0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test adding track
	err = conn.AddTrack(track)
	if err != nil {
		t.Fatalf("Failed to add track: %v", err)
	}

	// Test getting tracks
	tracks, err := conn.GetTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(tracks))
	}

	if tracks[0].Title != "Test Song" {
		t.Errorf("Expected track title to be 'Test Song', got '%s'", tracks[0].Title)
	}

	// Test getting track by ID
	retrievedTrack, err := conn.GetTrackByID(tracks[0].ID)
	if err != nil {
		t.Fatalf("Failed to get track by ID: %v", err)
	}

	if retrievedTrack.Title != "Test Song" {
		t.Errorf("Expected retrieved track title to be 'Test Song', got '%s'", retrievedTrack.Title)
	}
}

func createTestAudioFiles(t *testing.T, dir string) []string {
	// Create test audio files (mock files for testing)
	testFiles := []string{
		"test1.mp3",
		"test2.mp3",
		"test3.mp3",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(dir, filename)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
	}

	return testFiles
}

// createTestTracksDirectly creates test tracks directly in the database
func createTestTracksDirectly(t *testing.T, app *app.App) {
	// Create test tracks directly in the database
	testTracks := []database.PersistentTrack{
		{
			FilePath:        "/tmp/test1.mp3",
			FileName:        "test1.mp3",
			Title:           "Test Song 1",
			AlbumName:       "Test Album",
			ArtistName:      "Test Artist",
			AlbumArtistName: stringPtr("Test Artist"),
			Duration:        180.0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			FilePath:        "/tmp/test2.mp3",
			FileName:        "test2.mp3",
			Title:           "Test Song 2",
			AlbumName:       "Test Album",
			ArtistName:      "Test Artist",
			AlbumArtistName: stringPtr("Test Artist"),
			Duration:        200.0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			FilePath:        "/tmp/test3.mp3",
			FileName:        "test3.mp3",
			Title:           "Test Song 3",
			AlbumName:       "Test Album 2",
			ArtistName:      "Test Artist 2",
			AlbumArtistName: stringPtr("Test Artist 2"),
			Duration:        220.0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	// Add tracks directly to database
	for _, track := range testTracks {
		err := app.AddTrack(&track)
		if err != nil {
			t.Fatalf("Failed to add test track: %v", err)
		}
	}
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

func TestConcurrentOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create database connection
	conn, err := database.NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer conn.Close()

	// Test concurrent track additions
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			track := &database.PersistentTrack{
				FilePath:   "/tmp/test" + string(rune(id)) + ".mp3",
				FileName:   "test" + string(rune(id)) + ".mp3",
				Title:      "Test Song " + string(rune(id)),
				AlbumName:  "Test Album",
				ArtistName: "Test Artist",
				Duration:   180.0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			err := conn.AddTrack(track)
			if err != nil {
				t.Errorf("Failed to add track %d: %v", id, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test that all tracks were added
	tracks, err := conn.GetTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 10 {
		t.Errorf("Expected 10 tracks, got %d", len(tracks))
	}
}

func TestErrorHandling(t *testing.T) {
	// Test with invalid directory
	_, err := database.NewConnection("/invalid/path/that/does/not/exist")
	if err == nil {
		t.Errorf("Expected error for invalid directory, got nil")
	}

	// Test with empty directory
	tempDir := t.TempDir()
	conn, err := database.NewConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer conn.Close()

	// Test getting non-existent track
	_, err = conn.GetTrackByID(999999)
	if err == nil {
		t.Errorf("Expected error for non-existent track, got nil")
	}
}
