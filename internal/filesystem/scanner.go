package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"lrcget-go/internal/database"

	"github.com/dhowden/tag"
)

// Scanner represents a file system scanner
type Scanner struct {
}

// NewScanner creates a new file system scanner
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanDirectories scans directories for audio files
func (s *Scanner) ScanDirectories(directories []string) ([]database.PersistentTrack, error) {
	var allTracks []database.PersistentTrack

	for _, directory := range directories {
		tracks, err := s.scanDirectory(directory)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %w", directory, err)
		}
		allTracks = append(allTracks, tracks...)
	}

	return allTracks, nil
}

// scanDirectory scans a single directory for audio files
func (s *Scanner) scanDirectory(directory string) ([]database.PersistentTrack, error) {
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

		// Extract metadata from audio file
		track, err := s.extractMetadata(path)
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

// isAudioFile checks if a file is an audio file
func (s *Scanner) isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	audioExtensions := []string{".mp3", ".m4a", ".flac", ".ogg", ".opus", ".wav", ".aac", ".wma"}

	for _, audioExt := range audioExtensions {
		if ext == audioExt {
			return true
		}
	}

	return false
}

// extractMetadata extracts metadata from an audio file
func (s *Scanner) extractMetadata(filePath string) (*database.PersistentTrack, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read tags from file
	tags, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read tags: %w", err)
	}

	// Extract basic metadata
	title := tags.Title()
	if title == "" {
		// Use filename without extension as title
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	album := tags.Album()
	if album == "" {
		album = "Unknown Album"
	}

	artist := tags.Artist()
	if artist == "" {
		artist = "Unknown Artist"
	}

	albumArtist := tags.AlbumArtist()
	if albumArtist == "" {
		albumArtist = artist
	}

	// Get duration (in seconds) - for now, we'll set it to 0
	// In a real implementation, you would need to read the audio file properties
	var duration float64 = 0

	// Get track number
	trackNumber, _ := tags.Track()

	// Check for existing lyrics files
	txtLyrics := s.getTxtLyrics(filePath)
	lrcLyrics := s.getLrcLyrics(filePath)

	// Create track
	track := &database.PersistentTrack{
		FilePath:        filePath,
		FileName:        filepath.Base(filePath),
		Title:           title,
		AlbumName:       album,
		AlbumArtistName: &albumArtist,
		ArtistName:      artist,
		Duration:        duration,
		Instrumental:    false,
		TxtLyrics:       txtLyrics,
		LrcLyrics:       lrcLyrics,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if trackNumber > 0 {
		trackNum := int64(trackNumber)
		track.TrackNumber = &trackNum
	}

	// Set lowercase versions for searching
	titleLower := strings.ToLower(title)
	track.TitleLower = &titleLower

	return track, nil
}

// getTxtLyrics looks for a .txt lyrics file
func (s *Scanner) getTxtLyrics(filePath string) *string {
	txtPath := s.getTxtPath(filePath)
	content, err := os.ReadFile(txtPath)
	if err != nil {
		return nil
	}

	lyrics := string(content)
	return &lyrics
}

// getLrcLyrics looks for a .lrc lyrics file
func (s *Scanner) getLrcLyrics(filePath string) *string {
	lrcPath := s.getLrcPath(filePath)
	content, err := os.ReadFile(lrcPath)
	if err != nil {
		return nil
	}

	lyrics := string(content)
	return &lyrics
}

// getTxtPath returns the path to the .txt lyrics file
func (s *Scanner) getTxtPath(filePath string) string {
	dir := filepath.Dir(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return filepath.Join(dir, base+".txt")
}

// getLrcPath returns the path to the .lrc lyrics file
func (s *Scanner) getLrcPath(filePath string) string {
	dir := filepath.Dir(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return filepath.Join(dir, base+".lrc")
}

// CountFiles counts the number of audio files in directories
func (s *Scanner) CountFiles(directories []string) (int, error) {
	count := 0

	for _, directory := range directories {
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if s.isAudioFile(path) {
				count++
			}

			return nil
		})

		if err != nil {
			return 0, fmt.Errorf("failed to count files in directory %s: %w", directory, err)
		}
	}

	return count, nil
}
