package audio

import (
	"sync"
	"time"

	"lrcget-go/internal/database"
)

// PlayerStatus represents the current player status
type PlayerStatus int

const (
	Stopped PlayerStatus = iota
	Playing
	Paused
)

// PlayerState represents the current state of the audio player
type PlayerState struct {
	Status   PlayerStatus `json:"status"`
	Progress float64      `json:"progress"`
	Duration float64      `json:"duration"`
	Volume   float64      `json:"volume"`
	Track    *database.PersistentTrack `json:"track,omitempty"`
}

// Player represents an audio player
type Player struct {
	mu       sync.RWMutex
	status   PlayerStatus
	progress float64
	duration float64
	volume   float64
	track    *database.PersistentTrack
	startTime time.Time
	pausedTime time.Duration
}

// NewPlayer creates a new audio player
func NewPlayer() (*Player, error) {
	// Note: In a real implementation, you would initialize the audio system here
	// For now, we'll create a mock player that tracks state
	
	return &Player{
		status:   Stopped,
		progress: 0,
		duration: 0,
		volume:   1.0,
	}, nil
}

// Play starts playing a track
func (p *Player) Play(track *database.PersistentTrack) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Stop current playback if any
	p.status = Stopped
	
	// Set new track
	p.track = track
	p.duration = track.Duration
	p.progress = 0
	p.startTime = time.Now()
	p.pausedTime = 0
	p.status = Playing
	
	// In a real implementation, you would start audio playback here
	// For now, we'll just update the state
	
	return nil
}

// Pause pauses the current track
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.status == Playing {
		p.pausedTime += time.Since(p.startTime)
		p.status = Paused
	}
}

// Resume resumes the paused track
func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.status == Paused {
		p.startTime = time.Now()
		p.status = Playing
	}
}

// Stop stops the current track
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.status = Stopped
	p.progress = 0
	p.pausedTime = 0
}

// Seek seeks to a specific position in the track
func (p *Player) Seek(position float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if position < 0 {
		position = 0
	}
	if position > p.duration {
		position = p.duration
	}
	
	p.progress = position
	p.startTime = time.Now().Add(-time.Duration(position) * time.Second)
	p.pausedTime = 0
}

// SetVolume sets the player volume
func (p *Player) SetVolume(volume float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	
	p.volume = volume
}

// GetState returns the current player state
func (p *Player) GetState() PlayerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Calculate current progress
	var currentProgress float64
	if p.status == Playing {
		elapsed := time.Since(p.startTime) + p.pausedTime
		currentProgress = elapsed.Seconds()
		if currentProgress > p.duration {
			currentProgress = p.duration
		}
	} else {
		currentProgress = p.progress
	}
	
	return PlayerState{
		Status:   p.status,
		Progress: currentProgress,
		Duration: p.duration,
		Volume:   p.volume,
		Track:    p.track,
	}
}

// UpdateState updates the player state (called periodically)
func (p *Player) UpdateState() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.status == Playing {
		elapsed := time.Since(p.startTime) + p.pausedTime
		p.progress = elapsed.Seconds()
		
		// Check if track has finished
		if p.progress >= p.duration {
			p.status = Stopped
			p.progress = 0
		}
	}
}

// IsPlaying returns true if the player is currently playing
func (p *Player) IsPlaying() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.status == Playing
}

// IsPaused returns true if the player is currently paused
func (p *Player) IsPaused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.status == Paused
}

// IsStopped returns true if the player is stopped
func (p *Player) IsStopped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.status == Stopped
}

// GetCurrentTrack returns the currently loaded track
func (p *Player) GetCurrentTrack() *database.PersistentTrack {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.track
}

// GetProgress returns the current playback progress in seconds
func (p *Player) GetProgress() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.progress
}

// GetDuration returns the duration of the current track
func (p *Player) GetDuration() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.duration
}

// GetVolume returns the current volume
func (p *Player) GetVolume() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.volume
}
