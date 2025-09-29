package utils

import (
	"errors"
	"log"
	"strings"
)

// SafeError wraps internal errors with user-safe messages
type SafeError struct {
	Message string
	Err     error
}

func (e *SafeError) Error() string {
	return e.Message
}

func (e *SafeError) Unwrap() error {
	return e.Err
}

// HandleError logs detailed error internally and returns safe error to user
func HandleError(operation string, err error) error {
	log.Printf("Error in %s: %v", operation, err)
	return &SafeError{
		Message: "An error occurred. Please try again.",
		Err:     err,
	}
}

// HandleErrorWithMessage logs detailed error and returns custom safe message
func HandleErrorWithMessage(operation string, err error, userMessage string) error {
	log.Printf("Error in %s: %v", operation, err)
	return &SafeError{
		Message: userMessage,
		Err:     err,
	}
}

// HandleDatabaseError handles database-specific errors safely
func HandleDatabaseError(operation string, err error) error {
	log.Printf("Database error in %s: %v", operation, err)

	// Provide more specific messages for common database errors
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "no such table") {
			return &SafeError{
				Message: "Database not properly initialized. Please restart the application.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "database is locked") {
			return &SafeError{
				Message: "Database is currently in use. Please try again in a moment.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "disk I/O error") {
			return &SafeError{
				Message: "Storage error occurred. Please check your disk space and permissions.",
				Err:     err,
			}
		}
	}

	return &SafeError{
		Message: "Database operation failed. Please try again.",
		Err:     err,
	}
}

// HandleNetworkError handles network-related errors safely
func HandleNetworkError(operation string, err error) error {
	log.Printf("Network error in %s: %v", operation, err)

	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "timeout") {
			return &SafeError{
				Message: "Request timed out. Please check your internet connection and try again.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "connection refused") {
			return &SafeError{
				Message: "Unable to connect to server. Please check your internet connection.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "no such host") {
			return &SafeError{
				Message: "Server not found. Please check your internet connection.",
				Err:     err,
			}
		}
	}

	return &SafeError{
		Message: "Network error occurred. Please check your internet connection and try again.",
		Err:     err,
	}
}

// HandleFileError handles file system errors safely
func HandleFileError(operation string, err error) error {
	log.Printf("File error in %s: %v", operation, err)

	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "permission denied") {
			return &SafeError{
				Message: "Permission denied. Please check file permissions and try again.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "no such file") {
			return &SafeError{
				Message: "File not found. Please check the file path and try again.",
				Err:     err,
			}
		}
		if strings.Contains(errStr, "disk full") {
			return &SafeError{
				Message: "Insufficient disk space. Please free up space and try again.",
				Err:     err,
			}
		}
	}

	return &SafeError{
		Message: "File operation failed. Please check the file and try again.",
		Err:     err,
	}
}

// IsSafeError checks if an error is a SafeError
func IsSafeError(err error) bool {
	var safeErr *SafeError
	return errors.As(err, &safeErr)
}

// GetInternalError extracts the internal error from SafeError
func GetInternalError(err error) error {
	var safeErr *SafeError
	if errors.As(err, &safeErr) {
		return safeErr.Err
	}
	return err
}
