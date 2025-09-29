package interfaces

import (
	"lrcget-go/internal/audio"
	"lrcget-go/internal/database"
)

// AudioPlayerInterface defines the interface for audio player operations
type AudioPlayerInterface interface {
	// Playback control
	Play(track *database.PersistentTrack) error
	Pause() error
	Resume() error
	Stop() error
	Seek(position float64) error
	SetVolume(volume float64) error
	
	// State management
	GetState() audio.PlayerState
	GetPosition() float64
	GetDuration() float64
	GetVolume() float64
	IsPlaying() bool
	IsPaused() bool
	IsStopped() bool
	
	// Event handling
	OnStateChange(callback func(audio.PlayerState))
	OnPositionChange(callback func(float64))
	OnVolumeChange(callback func(float64))
	
	// Cleanup
	Close() error
}

// AudioDecoderInterface defines the interface for audio decoding
type AudioDecoderInterface interface {
	Decode(filePath string) (AudioData, error)
	GetDuration(filePath string) (float64, error)
	GetMetadata(filePath string) (AudioMetadata, error)
	Close() error
}

// AudioData represents decoded audio data
type AudioData interface {
	GetSamples() []float32
	GetSampleRate() int
	GetChannels() int
	GetDuration() float64
}

// AudioMetadata represents audio file metadata
type AudioMetadata interface {
	GetTitle() string
	GetArtist() string
	GetAlbum() string
	GetDuration() float64
	GetBitrate() int
	GetSampleRate() int
	GetChannels() int
}

// AudioStreamInterface defines the interface for audio streaming
type AudioStreamInterface interface {
	Start() error
	Stop() error
	Pause() error
	Resume() error
	Seek(position float64) error
	SetVolume(volume float64) error
	GetState() audio.PlayerState
	GetPosition() float64
	GetDuration() float64
	Close() error
}
