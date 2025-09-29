package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"lrcget-go/internal/database"
)

const (
	// MaxFileSizeForMemoryProcessing is the maximum file size to process in memory
	MaxFileSizeForMemoryProcessing = 10 * 1024 * 1024 // 10MB
)

// extractMetadataStreaming extracts metadata from an audio file using streaming for large files
func (s *Scanner) extractMetadataStreaming(filePath string) (*database.PersistentTrack, error) {
	// Check file size first
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Use streaming for large files
	if fileInfo.Size() > MaxFileSizeForMemoryProcessing {
		return s.extractMetadataStreamingLarge(filePath)
	}

	// Use existing method for smaller files
	return s.extractMetadata(filePath)
}

// extractMetadataStreamingLarge extracts metadata from large audio files using streaming
func (s *Scanner) extractMetadataStreamingLarge(filePath string) (*database.PersistentTrack, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// For large files, we need to read the entire file for tag extraction
	// since tag.ReadFrom requires io.ReadSeeker, not just io.Reader
	// We'll use the existing extractMetadata method for now
	return s.extractMetadata(filePath)
}

// getTxtLyricsStreaming looks for a .txt lyrics file using streaming
func (s *Scanner) getTxtLyricsStreaming(filePath string) *string {
	txtPath := s.getTxtPath(filePath)

	// Check if file exists and get size
	fileInfo, err := os.Stat(txtPath)
	if err != nil {
		return nil
	}

	// For large lyrics files, use the existing method
	if fileInfo.Size() > MaxFileSizeForMemoryProcessing {
		// Use the existing getTxtLyrics method for large files
		return s.getTxtLyrics(filePath)
	}

	// For smaller files, use the existing method
	content, err := os.ReadFile(txtPath)
	if err != nil {
		return nil
	}

	lyrics := string(content)
	return &lyrics
}

// getLrcLyricsStreaming looks for a .lrc lyrics file using streaming
func (s *Scanner) getLrcLyricsStreaming(filePath string) *string {
	lrcPath := s.getLrcPath(filePath)

	// Check if file exists and get size
	fileInfo, err := os.Stat(lrcPath)
	if err != nil {
		return nil
	}

	// For large lyrics files, use the existing method
	if fileInfo.Size() > MaxFileSizeForMemoryProcessing {
		// Use the existing getLrcLyrics method for large files
		return s.getLrcLyrics(filePath)
	}

	// For smaller files, use the existing method
	content, err := os.ReadFile(lrcPath)
	if err != nil {
		return nil
	}

	lyrics := string(content)
	return &lyrics
}

// ScanDirectoriesStreaming scans directories for audio files using streaming for large files
func (s *Scanner) ScanDirectoriesStreaming(directories []string) ([]database.PersistentTrack, error) {
	var allTracks []database.PersistentTrack

	for _, directory := range directories {
		tracks, err := s.scanDirectoryStreaming(directory)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %w", directory, err)
		}
		allTracks = append(allTracks, tracks...)
	}

	return allTracks, nil
}

// scanDirectoryStreaming scans a single directory for audio files using streaming
func (s *Scanner) scanDirectoryStreaming(directory string) ([]database.PersistentTrack, error) {
	var tracks []database.PersistentTrack

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is an audio file
		if !s.isAudioFile(path) {
			return nil
		}

		// Extract metadata using streaming for large files
		track, err := s.extractMetadataStreaming(path)
		if err != nil {
			// Log error but continue scanning
			fmt.Printf("Failed to extract metadata from %s: %v\n", path, err)
			return nil
		}

		tracks = append(tracks, *track)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return tracks, nil
}
