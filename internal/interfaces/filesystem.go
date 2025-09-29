package interfaces

import (
	"lrcget-go/internal/database"
)

// FileSystemInterface defines the interface for file system operations
type FileSystemInterface interface {
	// Directory scanning
	ScanDirectories(directories []string) ([]database.PersistentTrack, error)
	ScanDirectory(directory string) ([]database.PersistentTrack, error)
	CountFiles(directories []string) (int, error)
	
	// File operations
	IsAudioFile(path string) bool
	ExtractMetadata(filePath string) (*database.PersistentTrack, error)
	GetTxtLyrics(filePath string) *string
	GetLrcLyrics(filePath string) *string
	
	// Streaming operations
	ScanDirectoriesStreaming(directories []string) ([]database.PersistentTrack, error)
	ExtractMetadataStreaming(filePath string) (*database.PersistentTrack, error)
	GetTxtLyricsStreaming(filePath string) *string
	GetLrcLyricsStreaming(filePath string) *string
	
	// Utility operations
	GetTxtPath(filePath string) string
	GetLrcPath(filePath string) string
}

// FileWatcherInterface defines the interface for file system watching
type FileWatcherInterface interface {
	Watch(directories []string) error
	Stop() error
	OnFileAdded(callback func(string))
	OnFileRemoved(callback func(string))
	OnFileModified(callback func(string))
	OnDirectoryAdded(callback func(string))
	OnDirectoryRemoved(callback func(string))
}

// FileCacheInterface defines the interface for file caching
type FileCacheInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}) error
	Delete(key string) error
	Clear() error
	Size() int
	Keys() []string
}

// FileValidatorInterface defines the interface for file validation
type FileValidatorInterface interface {
	ValidatePath(path string) error
	ValidateAudioFile(path string) error
	ValidateLyricsFile(path string) error
	IsValidAudioFormat(path string) bool
	IsValidLyricsFormat(path string) bool
}
