package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string
	Format     string
	Output     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// DefaultLogConfig returns the default logging configuration
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}
}

// InitLogger initializes the logger with the given configuration
func InitLogger(config *LogConfig) error {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	// Set output
	switch config.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		// File output
		if err := setupFileOutput(config.Output); err != nil {
			return err
		}
	}

	return nil
}

// setupFileOutput sets up file output for logging
func setupFileOutput(logFile string) error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger.SetOutput(file)
	return nil
}

// GetLogger returns the logger instance
func GetLogger() *log.Logger {
	if logger == nil {
		// Initialize with default config if not already initialized
		config := DefaultLogConfig()
		InitLogger(config)
	}
	return logger
}

// LogError logs an error with context
func LogError(operation string, err error) {
	GetLogger().Printf("ERROR [%s]: %v", operation, err)
}

// LogWarning logs a warning with context
func LogWarning(operation string, message string) {
	GetLogger().Printf("WARN [%s]: %s", operation, message)
}

// LogInfo logs an info message with context
func LogInfo(operation string, message string) {
	GetLogger().Printf("INFO [%s]: %s", operation, message)
}

// LogDebug logs a debug message with context
func LogDebug(operation string, message string) {
	GetLogger().Printf("DEBUG [%s]: %s", operation, message)
}

// LogPerformance logs performance metrics
func LogPerformance(operation string, duration time.Duration) {
	GetLogger().Printf("PERF [%s]: %v", operation, duration)
}

// LogDatabaseOperation logs database operations
func LogDatabaseOperation(operation string, table string, duration time.Duration, err error) {
	if err != nil {
		GetLogger().Printf("DB_ERROR [%s:%s]: %v (duration: %v)", operation, table, err, duration)
	} else {
		GetLogger().Printf("DB_SUCCESS [%s:%s]: completed in %v", operation, table, duration)
	}
}

// LogNetworkOperation logs network operations
func LogNetworkOperation(operation string, url string, duration time.Duration, statusCode int, err error) {
	if err != nil {
		GetLogger().Printf("NET_ERROR [%s:%s]: %v (status: %d, duration: %v)", operation, url, err, statusCode, duration)
	} else {
		GetLogger().Printf("NET_SUCCESS [%s:%s]: status %d in %v", operation, url, statusCode, duration)
	}
}

// LogFileOperation logs file operations
func LogFileOperation(operation string, filePath string, duration time.Duration, err error) {
	if err != nil {
		GetLogger().Printf("FILE_ERROR [%s:%s]: %v (duration: %v)", operation, filePath, err, duration)
	} else {
		GetLogger().Printf("FILE_SUCCESS [%s:%s]: completed in %v", operation, filePath, duration)
	}
}

// LogSecurityEvent logs security-related events
func LogSecurityEvent(event string, details string) {
	GetLogger().Printf("SECURITY [%s]: %s", event, details)
}

// LogAuditEvent logs audit events
func LogAuditEvent(action string, resource string, userID string, details string) {
	GetLogger().Printf("AUDIT [%s:%s:%s]: %s", action, resource, userID, details)
}

// LogStartup logs application startup
func LogStartup(version string, config string) {
	GetLogger().Printf("STARTUP [%s]: %s", version, config)
}

// LogShutdown logs application shutdown
func LogShutdown(reason string) {
	GetLogger().Printf("SHUTDOWN: %s", reason)
}
