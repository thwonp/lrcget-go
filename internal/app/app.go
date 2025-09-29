package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"lrcget-go/internal/audio"
	"lrcget-go/internal/database"
	"lrcget-go/internal/filesystem"
	"lrcget-go/internal/lrclib"
)

// App represents the main application
type App struct {
	ctx     context.Context
	db      *database.Connection
	player  *audio.Player
	scanner *filesystem.Scanner
	lrclib  *lrclib.Client
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the application starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	dataDir := a.getDataDirectory()
	db, err := database.NewConnection(dataDir)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		// Create a dummy connection to prevent nil pointer dereference
		// In a real application, this should be handled differently
		return
	}
	a.db = db

	// Initialize audio player
	player, err := audio.NewPlayer()
	if err != nil {
		fmt.Printf("Failed to initialize audio player: %v\n", err)
		return
	}
	a.player = player

	// Initialize scanner
	a.scanner = filesystem.NewScanner()

	// Initialize LRCLIB client
	config, err := a.db.GetConfig()
	if err != nil {
		config = &database.PersistentConfig{
			LrclibInstance: "https://lrclib.net",
		}
	}
	a.lrclib = lrclib.NewClient(config.LrclibInstance)

	// Start background tasks
	go a.startAudioStateUpdater()
}

// OnDomReady is called when the DOM is ready
func (a *App) OnDomReady(ctx context.Context) {
	// Frontend is ready
}

// OnShutdown is called when the application shuts down
func (a *App) OnShutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// getDataDirectory returns the data directory path
func (a *App) getDataDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./data"
	}
	return filepath.Join(homeDir, ".lrcget")
}

// startAudioStateUpdater starts the audio state updater goroutine
func (a *App) startAudioStateUpdater() {
	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.player.UpdateState()
			// In a real implementation, you would emit the state to the frontend
		case <-a.ctx.Done():
			return
		}
	}
}
