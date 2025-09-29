package app

import (
	"fmt"

	"lrcget-go/internal/audio"
	"lrcget-go/internal/database"
	"lrcget-go/internal/lrclib"
	"lrcget-go/internal/utils"
)

// Directory management
func (a *App) GetDirectories() ([]string, error) {
	return a.db.GetDirectories()
}

func (a *App) SetDirectories(directories []string) error {
	// Validate each directory path
	for _, dir := range directories {
		if err := utils.ValidateDirectory(dir); err != nil {
			return utils.HandleErrorWithMessage("SetDirectories", err, "Invalid directory path provided")
		}
	}
	
	return a.db.SetDirectories(directories)
}

// Library management
func (a *App) GetInit() (bool, error) {
	return a.db.GetInit()
}

func (a *App) InitializeLibrary() error {
	go func() {
		// Get directories to scan
		directories, err := a.db.GetDirectories()
		if err != nil {
			fmt.Printf("Failed to get directories: %v\n", err)
			return
		}
		
		// Scan directories for tracks
		tracks, err := a.scanner.ScanDirectories(directories)
		if err != nil {
			fmt.Printf("Failed to scan directories: %v\n", err)
			return
		}
		
		// Add tracks to database
		for _, track := range tracks {
			err := a.db.AddTrack(&track)
			if err != nil {
				fmt.Printf("Failed to add track %s: %v\n", track.FilePath, err)
			}
		}
		
		// Mark library as initialized
		err = a.db.SetInit(true)
		if err != nil {
			fmt.Printf("Failed to set init status: %v\n", err)
		}
	}()
	
	return nil
}

// Track operations
func (a *App) GetTracks() ([]database.PersistentTrack, error) {
	return a.db.GetTracks()
}

func (a *App) GetTrack(trackID int64) (*database.PersistentTrack, error) {
	return a.db.GetTrackByID(trackID)
}

func (a *App) AddTrack(track *database.PersistentTrack) error {
	return a.db.AddTrack(track)
}

// Album operations
func (a *App) GetAlbums() ([]database.PersistentAlbum, error) {
	return a.db.GetAlbums()
}

func (a *App) GetAlbum(albumID int64) (*database.PersistentAlbum, error) {
	return a.db.GetAlbumByID(albumID)
}

func (a *App) GetTracksByAlbum(albumID int64) ([]database.PersistentTrack, error) {
	return a.db.GetTracksByAlbumID(albumID)
}

// Artist operations
func (a *App) GetArtists() ([]database.PersistentArtist, error) {
	return a.db.GetArtists()
}

func (a *App) GetArtist(artistID int64) (*database.PersistentArtist, error) {
	return a.db.GetArtistByID(artistID)
}

func (a *App) GetTracksByArtist(artistID int64) ([]database.PersistentTrack, error) {
	return a.db.GetTracksByArtistID(artistID)
}

// Lyrics operations
func (a *App) DownloadLyrics(trackID int64) (string, error) {
	track, err := a.db.GetTrackByID(trackID)
	if err != nil {
		return "", fmt.Errorf("failed to get track: %w", err)
	}
	
	config, err := a.db.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	
	// Update LRCLIB client with current instance
	a.lrclib.SetBaseURL(config.LrclibInstance)
	
	// Get lyrics from LRCLIB
	response, err := a.lrclib.GetLyrics(a.ctx, track.Title, track.AlbumName, track.ArtistName, track.Duration)
	if err != nil {
		return "", fmt.Errorf("failed to get lyrics: %w", err)
	}
	
	switch resp := response.(type) {
	case lrclib.SyncedLyrics:
		err = a.db.UpdateTrackSyncedLyrics(trackID, resp.Synced, resp.Plain)
		if err != nil {
			return "", fmt.Errorf("failed to update synced lyrics: %w", err)
		}
		return "Synced lyrics downloaded", nil
		
	case lrclib.UnsyncedLyrics:
		err = a.db.UpdateTrackPlainLyrics(trackID, resp.Plain)
		if err != nil {
			return "", fmt.Errorf("failed to update plain lyrics: %w", err)
		}
		return "Plain lyrics downloaded", nil
		
	case lrclib.Instrumental:
		err = a.db.UpdateTrackInstrumental(trackID)
		if err != nil {
			return "", fmt.Errorf("failed to update instrumental: %w", err)
		}
		return "Marked track as instrumental", nil
		
	case lrclib.None:
		return "", fmt.Errorf("lyrics not found")
	}
	
	return "", fmt.Errorf("unknown response type")
}

// Search lyrics
func (a *App) SearchLyrics(title, artist, album, query string) (*lrclib.SearchResponse, error) {
	// Validate and sanitize inputs
	title = utils.SanitizeInput(title)
	artist = utils.SanitizeInput(artist)
	album = utils.SanitizeInput(album)
	query = utils.SanitizeInput(query)
	
	// Validate search query
	if err := utils.ValidateSearchQuery(query); err != nil {
		return nil, utils.HandleErrorWithMessage("SearchLyrics", err, "Invalid search query")
	}
	
	config, err := a.db.GetConfig()
	if err != nil {
		return nil, utils.HandleDatabaseError("GetConfig", err)
	}
	
	// Validate LRCLIB URL
	if err := utils.ValidateURL(config.LrclibInstance); err != nil {
		return nil, utils.HandleErrorWithMessage("SearchLyrics", err, "Invalid LRCLIB URL configuration")
	}
	
	// Update LRCLIB client with current instance
	a.lrclib.SetBaseURL(config.LrclibInstance)
	
	response, err := a.lrclib.SearchLyrics(a.ctx, title, artist, album, query)
	if err != nil {
		return nil, utils.HandleNetworkError("SearchLyrics", err)
	}
	
	return response, nil
}

// Publish lyrics
func (a *App) PublishLyrics(title, artist, album string, duration float64, syncedLyrics, plainLyrics *string, instrumental bool) (*lrclib.PublishResponse, error) {
	// Validate and sanitize inputs
	title = utils.SanitizeInput(title)
	artist = utils.SanitizeInput(artist)
	album = utils.SanitizeInput(album)
	
	// Validate duration
	if duration <= 0 {
		return nil, utils.HandleErrorWithMessage("PublishLyrics", fmt.Errorf("invalid duration"), "Duration must be positive")
	}
	
	config, err := a.db.GetConfig()
	if err != nil {
		return nil, utils.HandleDatabaseError("GetConfig", err)
	}
	
	// Validate LRCLIB URL
	if err := utils.ValidateURL(config.LrclibInstance); err != nil {
		return nil, utils.HandleErrorWithMessage("PublishLyrics", err, "Invalid LRCLIB URL configuration")
	}
	
	// Update LRCLIB client with current instance
	a.lrclib.SetBaseURL(config.LrclibInstance)
	
	req := lrclib.PublishRequest{
		TrackName:    title,
		ArtistName:   artist,
		AlbumName:    album,
		Duration:     duration,
		SyncedLyrics: syncedLyrics,
		PlainLyrics:  plainLyrics,
		Instrumental: instrumental,
	}
	
	response, err := a.lrclib.PublishLyrics(a.ctx, req)
	if err != nil {
		return nil, utils.HandleNetworkError("PublishLyrics", err)
	}
	
	return response, nil
}

// Flag lyrics
func (a *App) FlagLyrics(trackID int64, reason string) error {
	// Validate track ID
	if trackID <= 0 {
		return utils.HandleErrorWithMessage("FlagLyrics", fmt.Errorf("invalid track ID"), "Invalid track ID")
	}
	
	// Sanitize reason
	reason = utils.SanitizeInput(reason)
	
	config, err := a.db.GetConfig()
	if err != nil {
		return utils.HandleDatabaseError("GetConfig", err)
	}
	
	// Validate LRCLIB URL
	if err := utils.ValidateURL(config.LrclibInstance); err != nil {
		return utils.HandleErrorWithMessage("FlagLyrics", err, "Invalid LRCLIB URL configuration")
	}
	
	// Update LRCLIB client with current instance
	a.lrclib.SetBaseURL(config.LrclibInstance)
	
	req := lrclib.FlagRequest{
		TrackID: trackID,
		Reason:  reason,
	}
	
	err = a.lrclib.FlagLyrics(a.ctx, req)
	if err != nil {
		return utils.HandleNetworkError("FlagLyrics", err)
	}
	
	return nil
}

// Audio operations
func (a *App) PlayTrack(trackID int64) error {
	track, err := a.db.GetTrackByID(trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}
	
	return a.player.Play(track)
}

func (a *App) PauseTrack() error {
	a.player.Pause()
	return nil
}

func (a *App) ResumeTrack() error {
	a.player.Resume()
	return nil
}

func (a *App) StopTrack() error {
	a.player.Stop()
	return nil
}

func (a *App) SeekTrack(position float64) error {
	a.player.Seek(position)
	return nil
}

func (a *App) SetVolume(volume float64) error {
	a.player.SetVolume(volume)
	return nil
}

func (a *App) GetPlayerState() audio.PlayerState {
	return a.player.GetState()
}

// Configuration operations
func (a *App) GetConfig() (*database.PersistentConfig, error) {
	return a.db.GetConfig()
}

func (a *App) UpdateConfig(config *database.PersistentConfig) error {
	return a.db.UpdateConfig(config)
}
