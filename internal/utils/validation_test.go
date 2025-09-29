package utils

import (
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected error
	}{
		{
			name:     "valid absolute path",
			path:     "/home/user/music/song.mp3",
			expected: nil,
		},
		{
			name:     "empty path",
			path:     "",
			expected: ErrEmptyPath,
		},
		{
			name:     "path with traversal",
			path:     "/home/user/../etc/passwd",
			expected: ErrPathTraversal,
		},
		{
			name:     "path with home directory",
			path:     "/home/user/~/secret",
			expected: ErrHomeDirectoryRef,
		},
		{
			name:     "path with suspicious pattern",
			path:     "/home/user//music",
			expected: ErrSuspiciousPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("ValidateFilePath() error = %v, expected %v", err, tt.expected)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("ValidateFilePath() error = %v, expected %v", err, tt.expected)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "input with HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "input with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "input with control characters",
			input:    "Hello\x01World\x02Test",
			expected: "HelloWorldTest",
		},
		{
			name:     "input with whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected error
	}{
		{
			name:     "valid query",
			query:    "Hello World",
			expected: nil,
		},
		{
			name:     "empty query",
			query:    "",
			expected: ErrEmptyQuery,
		},
		{
			name:     "query too long",
			query:    string(make([]byte, 1001)),
			expected: ErrQueryTooLong,
		},
		{
			name:     "query with SQL injection",
			query:    "'; DROP TABLE users; --",
			expected: ErrInvalidQuery,
		},
		{
			name:     "query with quotes",
			query:    "Hello \"World\"",
			expected: ErrInvalidQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSearchQuery(tt.query)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("ValidateSearchQuery() error = %v, expected %v", err, tt.expected)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("ValidateSearchQuery() error = %v, expected %v", err, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected error
	}{
		{
			name:     "valid HTTP URL",
			url:      "http://example.com",
			expected: nil,
		},
		{
			name:     "valid HTTPS URL",
			url:      "https://example.com",
			expected: nil,
		},
		{
			name:     "empty URL",
			url:      "",
			expected: ErrEmptyURL,
		},
		{
			name:     "URL without protocol",
			url:      "example.com",
			expected: ErrInvalidURL,
		},
		{
			name:     "URL too long",
			url:      "https://example.com/" + string(make([]byte, 2048)),
			expected: ErrURLTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("ValidateURL() error = %v, expected %v", err, tt.expected)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("ValidateURL() error = %v, expected %v", err, tt.expected)
			}
		})
	}
}

func TestValidateDirectory(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected error
	}{
		{
			name:     "valid directory",
			path:     "/home/user/music",
			expected: nil,
		},
		{
			name:     "invalid path",
			path:     "/home/user/../etc",
			expected: ErrPathTraversal,
		},
		{
			name:     "empty path",
			path:     "",
			expected: ErrEmptyPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectory(tt.path)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("ValidateDirectory() error = %v, expected %v", err, tt.expected)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("ValidateDirectory() error = %v, expected %v", err, tt.expected)
			}
		})
	}
}
