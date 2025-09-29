package utils

import (
	"errors"
	"html"
	"path/filepath"
	"regexp"
	"strings"
)

// Error constants
var (
	ErrEmptyPath        = errors.New("path cannot be empty")
	ErrPathTraversal    = errors.New("path traversal not allowed")
	ErrHomeDirectoryRef = errors.New("home directory references not allowed")
	ErrSuspiciousPath   = errors.New("suspicious path pattern detected")
	ErrEmptyQuery       = errors.New("search query cannot be empty")
	ErrQueryTooLong     = errors.New("search query too long")
	ErrInvalidQuery     = errors.New("invalid characters in search query")
	ErrEmptyURL         = errors.New("URL cannot be empty")
	ErrInvalidURL       = errors.New("URL must start with http:// or https://")
	ErrURLTooLong       = errors.New("URL too long")
)

// ValidateFilePath validates file paths to prevent path traversal attacks
func ValidateFilePath(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return errors.New("path traversal not allowed")
	}

	// Check for absolute path (recommended for security)
	if !filepath.IsAbs(path) {
		// Convert to absolute path for validation
		absPath, err := filepath.Abs(path)
		if err != nil {
			return errors.New("invalid path format")
		}
		path = absPath
	}

	// Additional security checks
	if strings.Contains(path, "~") {
		return errors.New("home directory references not allowed")
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"//",
		"\\\\",
		"/./",
		"\\.\\",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(path, pattern) {
			return errors.New("suspicious path pattern detected")
		}
	}

	return nil
}

// SanitizeInput sanitizes user input to prevent XSS and injection attacks
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// HTML escape to prevent XSS
	input = html.EscapeString(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newlines and tabs
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")

	return input
}

// ValidateSearchQuery validates search queries
func ValidateSearchQuery(query string) error {
	if len(query) == 0 {
		return errors.New("search query cannot be empty")
	}

	if len(query) > 1000 {
		return errors.New("search query too long")
	}

	// Check for SQL injection patterns
	sqlPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"union", "select", "insert", "update", "delete",
		"drop", "create", "alter", "exec", "execute",
	}

	queryLower := strings.ToLower(query)
	for _, pattern := range sqlPatterns {
		if strings.Contains(queryLower, pattern) {
			return errors.New("invalid characters in search query")
		}
	}

	return nil
}

// ValidateURL validates URL format
func ValidateURL(url string) error {
	if url == "" {
		return errors.New("URL cannot be empty")
	}

	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}

	if len(url) > 2048 {
		return errors.New("URL too long")
	}

	return nil
}

// ValidateDirectory validates directory paths
func ValidateDirectory(path string) error {
	if err := ValidateFilePath(path); err != nil {
		return err
	}

	// Additional directory-specific validation
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return errors.New("invalid directory path")
		}
		path = absPath
	}

	return nil
}
